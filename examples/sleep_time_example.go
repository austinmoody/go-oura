package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	threeDaysAgo := time.Now().Add(-72 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	sleepTimes, err := client.GetSleepTimes(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting Sleep Times: %v", err)
		return
	}

	if len(sleepTimes.Items) > 0 {
		fmt.Printf(
			"There were %d Sleep Times found for date range: %v - %v\n",
			len(sleepTimes.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Sleep Time ID: %s\n",
			sleepTimes.Items[0].ID,
		)

		sleepTime, err := client.GetSleepTime(sleepTimes.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting Sleep Time: %v", err)
			return
		}

		fmt.Printf("Single Sleep Time Recommendation: %s\n", sleepTime.Recommendation)

	} else {
		fmt.Printf(
			"No Sleep Times were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
