package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	spo2s, err := client.GetSpo2Readings(twoDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting Spo2 readings: %v", err)
		return
	}

	if len(spo2s.Items) > 0 {
		fmt.Printf(
			"There were %d Spo2 readings found for date range: %v - %v\n",
			len(spo2s.Items),
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Spo2 ID: %s\n",
			spo2s.Items[0].ID,
		)

		singleSpo2, err := client.GetSpo2Reading(spo2s.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single spo2 reading: %v", err)
			return
		}

		fmt.Printf("Single Spo2 Reading Average: %f\n", singleSpo2.Percentage.Average)

	} else {
		fmt.Printf(
			"No Spo2 readings were found for the date range: %v - %v",
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
