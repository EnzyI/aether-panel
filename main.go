package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func main() {
	// Khởi tạo kết nối với Docker Engine
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Kiểm tra thông tin hệ thống Docker
	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Printf("Lỗi kết nối Docker: %v\n", err)
		return
	}

	fmt.Printf("--- Aether Panel Daemon ---\n")
	fmt.Printf("Kết nối thành công với Docker!\n")
	fmt.Printf("Số lượng Container đang chạy: %d\n", info.ContainersRunning)
	fmt.Printf("Hệ điều hành: %s\n", info.OperatingSystem)
}
