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

	sleepDocs, err := client.GetSleeps(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting DailySleep Items: %v", err)
		return
	}

	if len(sleepDocs.Items) > 0 {
		fmt.Printf(
			"There were %d DailySleep Items found for date range: %v - %v\n",
			len(sleepDocs.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First DailySleeps ID: %s\n",
			sleepDocs.Items[0].ID,
		)

		singleSleepDoc, err := client.GetSleep(sleepDocs.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single sleep item: %v", err)
			return
		}

		// Match the field 'Score' in DailySleep struct or replace it with the correct field
		fmt.Printf("Single DailySleep Score: %d\n", singleSleepDoc.Score)

	} else {
		fmt.Printf(
			"No DailySleep Items were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
