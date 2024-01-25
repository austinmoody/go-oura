package main

import (
	"fmt"
	"os"
	"time"

	"github.com/austinmoody/go_oura"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	sleepDocs, err := client.GetSleeps(twoDaysAgo, oneDayAgo, nil)
	if err != nil {
		fmt.Printf("Error getting Sleep Items: %v", err)
		return
	}

	if len(sleepDocs.Items) > 0 {
		fmt.Printf(
			"There were %d Sleep Items found for date range: %v - %v\n",
			len(sleepDocs.Items),
			twoDaysAgo.Format("02-Jan-2006"),
			oneDayAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Sleeps ID: %s\n",
			sleepDocs.Items[0].ID,
		)

		singleSleepDoc, err := client.GetSleep(sleepDocs.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single sleep item: %v", err)
			return
		}

		fmt.Printf("Single Sleep Total Duration: %d\n", singleSleepDoc.TotalSleepDuration)

	} else {
		fmt.Printf(
			"No Sleep Items were found for the date range: %v - %v",
			twoDaysAgo.Format("02-Jan-2006"),
			oneDayAgo.Format("02-Jan-2006"),
		)
	}

}
