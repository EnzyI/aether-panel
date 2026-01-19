package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	fmt.Println("--- AETHER PANEL SYSTEM ---")
	ctx := context.Background()

	// Khoi tao Docker Client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Loi ket noi Docker: %v\n", err)
		return
	}
	fmt.Println("[1/3] Ket noi Docker thanh cong!")

	// Bat dau thiet lap Server Minecraft
	fmt.Println("[2/3] Dang kiem tra va thiet lap Minecraft Geyser Server...")
	err = SetupGeyserServer(ctx, cli)
	if err != nil {
		fmt.Printf("Loi thiet lap Server: %v\n", err)
		return
	}

	fmt.Println("[3/3] HE THONG DA SAN SANG!")
	fmt.Println("---------------------------")
	fmt.Println("IP: Dung IP cua Codespaces")
	fmt.Println("Port Java: 25565")
	fmt.Println("Port Bedrock (Mobile): 19132")
	fmt.Println("---------------------------")

	// Giu chuong trinh chay ngam
	for {
		time.Sleep(time.Hour)
	}
}

func SetupGeyserServer(ctx context.Context, cli *client.Client) error {
	const (
		imageName     = "itzg/minecraft-server:latest"
		containerName = "aether-mc-geyser"
	)

	// 1. Kiem tra xem container da ton tai chua
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				if c.State != "running" {
					fmt.Println("Dang khoi dong lai server...")
					return cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
				}
				fmt.Println("Server dang chay san roi!")
				return nil
			}
		}
	}

	// 2. Tai Image (Pull)
	fmt.Println("Dang tai Image Minecraft (co the mat vai phut)...")
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	// 3. Cau hinh Cong (Port)
	portBindings := nat.PortMap{
		"25565/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "25565"}},
		"19132/udp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "19132"}},
	}

	// 4. Cau hinh Container
	containerConfig := &container.Config{
		Image: imageName,
		Env: []string{
			"TYPE=PAPER",
			"EULA=TRUE",
			"VERSION=LATEST",
			"PLUGINS=geyser,floodgate",
		},
		ExposedPorts: nat.PortSet{"25565/tcp": struct{}{}, "19132/udp": struct{}{}},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}

	// 5. Tao va Chay
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}
