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

// Khai bao cau truc du lieu tra ve API
type ContainerView struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	State  string  `json:"state"`
	Memory float64 `json:"ram"`
}

type StatsResponse struct {
	Containers []ContainerView `json:"containers"`
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Loi Docker:", err)
		return
	}

	// 1. Route tra ve giao dien web
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	// 2. Route API lay thong tin container thuc te
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		res := StatsResponse{Containers: []ContainerView{}}
		for _, c := range list {
			mem := 0.0
			if c.State == "running" {
				stats, err := cli.ContainerStats(context.Background(), c.ID, false)
				if err == nil {
					var v types.StatsJSON
					json.NewDecoder(stats.Body).Decode(&v)
					mem = float64(v.MemoryStats.Usage) / 1024 / 1024
					stats.Body.Close()
				}
			}
			res.Containers = append(res.Containers, ContainerView{
				ID:    c.ID[:12],
				Name:  strings.TrimPrefix(c.Names[0], "/"),
				State: strings.Title(c.State),
				Memory: mem,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("🚀 AETHER PANEL ONLINE: http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

// Phan nay toi da sua lai cach noi chuoi JS de khong bi loi ky tu $ tren Go
var indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Aether Panel</title>
    <style>
        body { background: #06080f; color: white; font-family: sans-serif; padding: 30px; }
        .card { background: #0f121d; border: 1px solid #1e293b; padding: 20px; border-radius: 15px; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #1e293b; }
        .Running { color: #10b981; }
        .Stopped { color: #ef4444; }
    </style>
</head>
<body>
    <div class="card">
        <h1>⚡ AETHER PANEL v3.0</h1>
        <table>
            <thead><tr><th>ID</th><th>TÊN</th><th>TRẠNG THÁI</th><th>RAM (MB)</th></tr></thead>
            <tbody id="rows"></tbody>
        </table>
    </div>
    <script>
        async function refresh() {
            const r = await fetch('/api/stats');
            const d = await r.json();
            const b = document.getElementById('rows');
            b.innerHTML = d.containers.map(c => 
                "<tr><td>"+c.id+"</td><td>"+c.name+"</td><td class='"+c.state+"'>"+c.state+"</td><td>"+c.ram.toFixed(1)+"</td></tr>"
            ).join('');
        }
        setInterval(refresh, 3000); refresh();
    </script>
</body>
</html>`
