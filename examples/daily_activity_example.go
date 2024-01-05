package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	activity, err := client.GetActivity("4xxx5173cbe-ef26-430f-adc4-c4a1424b45ab")
	if err != nil {
		fmt.Printf("Error getting activity: %v", err)
		return
	}

	fmt.Printf("Activity ID: %s", activity.ID)
}
