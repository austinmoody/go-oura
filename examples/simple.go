package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
)

func main() {
	client := go_oura.NewClient("")

	fmt.Println(client)
}
