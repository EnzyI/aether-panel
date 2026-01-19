package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Data models cho API
type ContainerView struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	State  string  `json:"state"`
	CPU    float64 `json:"cpu"` // %
	Memory float64 `json:"ram"` // MB
}

type StatsResponse struct {
	UpdatedAt  string          `json:"updatedAt"`
	Containers []ContainerView `json:"containers"`
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Route tra ve Giao dien HTML cua bro
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	// Route API tra ve du lieu Container thuc te
	http.HandleFunc("/api/stats", apiStatsHandler(cli))

	fmt.Println("🚀 AETHER PANEL v3.0 dang chay tai: http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Ham tinh toan % CPU (Docker SDK standard)
func cpuPercent(prevCPU, prevSystem uint64, v *types.StatsJSON) float64 {
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - prevCPU)
	sysDelta := float64(v.CPUStats.SystemUsage - prevSystem)
	if sysDelta > 0 && cpuDelta > 0 {
		cores := float64(len(v.CPUStats.CPUUsage.PercpuUsage))
		if cores == 0 { cores = 1 }
		return (cpuDelta / sysDelta) * cores * 100.0
	}
	return 0.0
}

func apiStatsHandler(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		out := StatsResponse{
			UpdatedAt:  time.Now().Format(time.RFC3339),
			Containers: make([]ContainerView, 0, len(list)),
		}

		for _, c := range list {
			view := ContainerView{
				ID:    c.ID[:12],
				Name:  strings.TrimPrefix(c.Names[0], "/"),
				State: strings.Title(c.State),
			}

			if c.State == "running" {
				resp, err := cli.ContainerStats(ctx, c.ID, false)
				if err == nil {
					var v types.StatsJSON
					if json.NewDecoder(resp.Body).Decode(&v) == nil {
						view.Memory = float64(v.MemoryStats.Usage) / 1024 / 1024
						view.CPU = cpuPercent(v.PreCPUStats.CPUUsage.TotalUsage, v.PreCPUStats.SystemUsage, &v)
					}
					resp.Body.Close()
				}
			}
			out.Containers = append(out.Containers, view)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

// Giao dien HTML ket hop giua THIET KE cua bro va LOGIC cua Sybau
var indexHTML = `
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Aether | Game Control Panel</title>
    <script src="https://unpkg.com/lucide@latest"></script>
    <style>
        :root {
            --bg-dark: #06080f; --bg-card: #0f121d; --primary: #3b82f6;
            --text: #f8fafc; --text-muted: #64748b; --success: #10b981;
            --danger: #ef4444; --border: rgba(255, 255, 255, 0.08);
        }
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Inter', sans-serif; background: var(--bg-dark); color: var(--text); display: flex; height: 100vh; overflow: hidden; }
        aside { width: 260px; background: var(--bg-card); border-right: 1px solid var(--border); padding: 24px; display: flex; flex-direction: column; }
        .logo { font-size: 22px; font-weight: 800; color: var(--primary); margin-bottom: 40px; display: flex; align-items: center; gap: 10px; font-style: italic; }
        nav a { display: flex; align-items: center; gap: 12px; padding: 12px; color: var(--text-muted); text-decoration: none; border-radius: 8px; margin-bottom: 4px; cursor: pointer; transition: 0.2s; }
        nav a:hover, nav a.active { background: rgba(59, 130, 246, 0.1); color: var(--primary); }
        main { flex: 1; padding: 40px; overflow-y: auto; }
        .card { background: var(--bg-card); border-radius: 16px; border: 1px solid var(--border); padding: 20px; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th { text-align: left; color: var(--text-muted); font-size: 12px; padding: 12px; border-bottom: 1px solid var(--border); }
        td { padding: 14px; border-bottom: 1px solid var(--border); font-size: 14px; }
        .status-badge { padding: 4px 8px; border-radius: 4px; font-size: 11px; font-weight: bold; }
        .Running { background: rgba(16, 185, 129, 0.1); color: var(--success); }
        .Stopped { background: rgba(239, 68, 68, 0.1); color: var(--danger); }
    </style>
</head>
<body>
    <aside>
        <div class="logo"><i data-lucide="zap"></i> AETHER</div>
        <nav id="sidebar-nav">
            <a class="active"><i data-lucide="layout-dashboard"></i> Dashboard</a>
            <a><i data-lucide="files"></i> File Manager</a>
            <a><i data-lucide="database"></i> Databases</a>
            <a><i data-lucide="settings"></i> Settings</a>
        </nav>
    </aside>
    <main>
        <h1 style="margin-bottom: 25px;">Server Overview</h1>
        <div class="card">
            <h3 style="margin-bottom:15px">Live Container Stats</h3>
            <table>
                <thead>
                    <tr><th>ID</th><th>Container Name</th><th>Status</th><th>CPU %</th><th>RAM (MB)</th></tr>
                </thead>
                <tbody id="container-table">
                    </tbody>
            </table>
            <p style="font-size:10px; color:var(--text-muted); margin-top:15px">Auto-refresh: 3s | Aether Panel Core v3.0</p>
        </div>
    </main>

    <script>
        async function updateStats() {
            try {
                const response = await fetch('/api/stats');
                const data = await response.json();
                const tableBody = document.getElementById('container-table');
                tableBody.innerHTML = '';

                data.containers.forEach(c => {
                    tableBody.innerHTML += `
                        <tr>
                            <td>${c.id}</td>
                            <td><strong>${c.name}</strong></td>
                            <td><span class="status-badge ${c.state}">${c.state}</span></td>
                            <td>${c.cpu.toFixed(2)}%</td>
                            <td>${c.ram.toFixed(1)} MB</td>
                        </tr>
                    `;
                });
            } catch (err) { console.error("Loi update:", err); }
        }
        lucide.createIcons();
        updateStats();
        setInterval(updateStats, 3000);
    </script>
</body>
</html>
`
