package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Loi khoi tao: ", err)
		return
	}
	
	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Println("Loi ket noi Docker (Co the thieu quyen sudo): ", err)
		return
	}

	fmt.Println("--- AETHER PANEL ---")
	fmt.Printf("THANH CONG! Docker OS: %s\n", info.OperatingSystem)
}
