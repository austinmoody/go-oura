package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	threeDaysAgo := time.Now().Add(-72 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	ringConfigs, err := client.GetRingConfigurations(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting Ring Configurations: %v", err)
		return
	}

	if len(ringConfigs.Items) > 0 {
		fmt.Printf("There were %d Ring Configurations found\n", len(ringConfigs.Items))

		fmt.Printf("First Configuration ID: %s\n", ringConfigs.Items[0].ID)

		singleRingConfig, err := client.GetRingConfiguration(ringConfigs.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single ring configuration: %v", err)
			return
		}

		fmt.Printf("Single Configuration Color: %s\n", singleRingConfig.Color)
	} else {
		fmt.Println("No Ring Configurations were found.")
	}
}
