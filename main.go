package main

import (
	"fmt"
	"github.com/docker/docker/client"
)

func main() {
	_, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Loi ket noi Docker")
		return
	}
	fmt.Println("--- AETHER PANEL ---")
	fmt.Println("Ket noi Docker thanh cong!")
}
