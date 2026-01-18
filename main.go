package main

import (
	"fmt"
	"github.com/docker/docker/client"
)

func main() {
	_, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Loi: Khong ket noi duoc Docker")
		return
	}
	fmt.Println("--- AETHER PANEL ---")
	fmt.Println("THANH CONG: Da ket noi voi Docker Engine!")
}
