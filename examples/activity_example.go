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

	activities, err := client.GetActivities(threeDaysAgo, oneDaysAgo)
	if err != nil {
		fmt.Printf("Error getting activities: %v", err)
		return
	}

	if len(activities.Items) > 0 {
		fmt.Printf(
			"There were %d Activities found for date range: %v - %v\n",
			len(activities.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Activity ID: %s\n",
			activities.Items[0].ID,
		)

		singleActivity, err := client.GetActivity(activities.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single activity: %v", err)
			return
		}

		fmt.Printf("Single Activity Score: %d\n", singleActivity.Score)

	} else {
		fmt.Printf(
			"No activities were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}

}
