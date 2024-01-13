package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	personalInfo, err := client.GetPersonalInfo()
	if err != nil {
		fmt.Printf("Error getting Personal Info: %v", err)
		return
	}

	fmt.Printf("Personal Info: %+v\n", personalInfo)
}
