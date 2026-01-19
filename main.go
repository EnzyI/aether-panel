package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container" // Quan trọng cho v3.1
	"github.com/docker/docker/client"
)

// Data models
type ContainerView struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	State  string  `json:"state"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"ram"`
}

type StatsResponse struct {
	Containers []ContainerView `json:"containers"`
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	http.HandleFunc("/api/stats", apiStatsHandler(cli))
	http.HandleFunc("/api/action", apiActionHandler(cli))

	fmt.Println("🚀 AETHER PANEL v3.1 ONLINE: http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

// API Dieu khien Container
func apiActionHandler(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		cmd := r.URL.Query().Get("cmd")
		if id == "" || cmd == "" {
			http.Error(w, "Thieu ID hoac CMD", 400)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var err error
		switch cmd {
		case "start":
			err = cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
		case "stop":
			err = cli.ContainerStop(ctx, id, container.StopOptions{})
		case "restart":
			err = cli.ContainerRestart(ctx, id, container.StopOptions{})
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte(`{"status":"ok"}`))
	}
}

// API Lay Stats (CPU/RAM/List)
func apiStatsHandler(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		out := StatsResponse{Containers: []ContainerView{}}
		for _, c := range list {
			view := ContainerView{
				ID:    c.ID[:12],
				Name:  strings.TrimPrefix(c.Names[0], "/"),
				State: strings.Title(c.State),
			}
			// Lay RAM thuc te neu dang chay
			if c.State == "running" {
				stats, err := cli.ContainerStats(context.Background(), c.ID, false)
				if err == nil {
					var v types.StatsJSON
					json.NewDecoder(stats.Body).Decode(&v)
					view.Memory = float64(v.MemoryStats.Usage) / 1024 / 1024
					stats.Body.Close()
				}
			}
			out.Containers = append(out.Containers, view)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

// Giao dien HTML v3.1 cua Bro & Sybau
var indexHTML = `
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Aether Panel v3.1</title>
    <style>
        body { font-family: sans-serif; background: #06080f; color: #e5e7eb; padding: 20px; }
        .card { background: #0f121d; border: 1px solid #1f2937; border-radius: 12px; padding: 16px; max-width: 1000px; margin: auto; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #1f2937; }
        .badge { padding: 4px 8px; border-radius: 6px; font-size: 12px; font-weight: bold; }
        .Running { background: #064e3b; color: #a7f3d0; }
        .Exited { background: #3f3f46; color: #fff; }
        button { padding: 6px 12px; border-radius: 6px; border: none; cursor: pointer; font-size: 11px; margin-right: 4px; font-weight: bold; }
        .start { background: #16a34a; color: white; }
        .stop { background: #dc2626; color: white; }
        .restart { background: #f59e0b; color: black; }
    </style>
</head>
<body>
    <div class="card">
        <h1>⚡ AETHER PANEL v3.1</h1>
        <table>
            <thead>
                <tr>
                    <th>NAME</th><th>STATUS</th><th>RAM (MB)</th><th>ACTIONS</th>
                </tr>
            </thead>
            <tbody id="rows"></tbody>
        </table>
        <p style="font-size: 12px; color: #64748b; margin-top: 15px;">Auto refresh: 3s | Dev by Community</p>
    </div>
    <script>
        async function action(id, cmd) {
            console.log("Lenh: " + cmd);
            await fetch('/api/action?id=' + id + '&cmd=' + cmd);
            loadStats();
        }
        async function loadStats() {
            try {
                const res = await fetch('/api/stats');
                const data = await res.json();
                document.getElementById('rows').innerHTML = data.containers.map(c => 
                    "<tr>" +
                        "<td><strong>" + c.name + "</strong></td>" +
                        "<td><span class='badge " + c.state + "'>" + c.state + "</span></td>" +
                        "<td>" + c.ram.toFixed(1) + "</td>" +
                        "<td>" +
                            "<button class='start' onclick=\"action('"+c.id+"','start')\">Start</button>" +
                            "<button class='stop' onclick=\"action('"+c.id+"','stop')\">Stop</button>" +
                            "<button class='restart' onclick=\"action('"+c.id+"','restart')\">Restart</button>" +
                        "</td>" +
                    "</tr>"
                ).join('');
            } catch (e) {}
        }
        setInterval(loadStats, 3000); loadStats();
    </script>
</body>
</html>`
