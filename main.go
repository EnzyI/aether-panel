package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerView struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func main() {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		list, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		containers := []ContainerView{}
		for _, c := range list {
			containers = append(containers, ContainerView{
				ID: c.ID[:12], Name: strings.TrimPrefix(c.Names[0], "/"), State: strings.Title(c.State),
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"containers": containers})
	})

	http.HandleFunc("/api/action", func(w http.ResponseWriter, r *http.Request) {
		id, cmd := r.URL.Query().Get("id"), r.URL.Query().Get("cmd")
		ctx := context.Background()
		var err error
		if cmd == "start" { err = cli.ContainerStart(ctx, id, types.ContainerStartOptions{}) }
		if cmd == "stop" { err = cli.ContainerStop(ctx, id, container.StopOptions{}) }
		if cmd == "restart" { err = cli.ContainerRestart(ctx, id, container.StopOptions{}) }
		if err != nil { http.Error(w, err.Error(), 500); return }
		fmt.Fprint(w, "OK")
	})

	fmt.Println("🚀 PANEL v3.1.1 ONLINE")
	http.ListenAndServe(":8080", nil)
}

var indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Aether Control</title>
    <style>
        body { background: #06080f; color: white; font-family: sans-serif; padding: 15px; }
        .card { background: #0f121d; border: 1px solid #1e293b; padding: 20px; border-radius: 15px; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; }
        td { padding: 15px 5px; border-bottom: 1px solid #1e293b; font-size: 14px; }
        .btn { padding: 12px; border: none; border-radius: 8px; cursor: pointer; font-weight: bold; width: 100%; margin: 5px 0; display: block; }
        .start { background: #10b981; color: white; }
        .stop { background: #ef4444; color: white; }
        .restart { background: #f59e0b; color: black; }
        .status-Running { color: #10b981; font-weight: bold; }
    </style>
</head>
<body>
    <div class="card">
        <h1>⚡ AETHER v3.1.1</h1>
        <p id="msg" style="color: #64748b; font-size: 12px;">Đang sẵn sàng...</p>
        <table>
            <tbody id="rows"></tbody>
        </table>
    </div>
    <script>
        async function doAction(id, cmd) {
            document.getElementById('msg').innerText = "Đang " + cmd + "...";
            await fetch('/api/action?id=' + id + '&cmd=' + cmd);
            setTimeout(refresh, 1000);
        }
        async function refresh() {
            const r = await fetch('/api/stats');
            const d = await r.json();
            document.getElementById('rows').innerHTML = d.containers.map(c => 
                "<tr><td><strong>" + c.name + "</strong><br><span class='status-" + c.state + "'>" + c.state + "</span></td>" +
                "<td>" +
                    (c.state !== 'Running' ? "<button class='btn start' onclick=\"doAction('"+c.id+"','start')\">START</button>" : "") +
                    (c.state === 'Running' ? "<button class='btn stop' onclick=\"doAction('"+c.id+"','stop')\">STOP</button>" : "") +
                    "<button class='btn restart' onclick=\"doAction('"+c.id+"','restart')\">RESTART</button>" +
                "</td></tr>"
            ).join('');
        }
        setInterval(refresh, 3000); refresh();
    </script>
</body>
</html>`

