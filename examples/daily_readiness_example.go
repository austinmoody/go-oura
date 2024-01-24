package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"os"
	"time"
)

func main() {

	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	// Get Multiple DailyReadiness Items (if they exist) for the last couple days.
	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	readiness, err := client.GetReadinesses(twoDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting multiple readinesses: %v", err)
		return
	}

	if len(readiness.Items) > 0 {
		fmt.Printf(
			"There were %d DailyReadiness Items found for date range: %v - %v\n",
			len(readiness.Items),
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First DailyReadiness ID: %s\n",
			readiness.Items[0].Id,
		)

		// Could of course print whatever from readiness.Items[0]
		// but will pull directly to use single document function
		singleReadiness, err := client.GetReadiness(readiness.Items[0].Id)
		if err != nil {
			fmt.Printf("Error getting single readiness: %v", err)
			return
		}

		fmt.Printf("Single DailyReadiness Score: %d\n", singleReadiness.Score)

	} else {
		fmt.Printf(
			"No DailyReadiness Items were found for the date range: %v - %v",
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
