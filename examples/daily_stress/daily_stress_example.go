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

	stressDocs, err := client.GetStresses(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting DailyStress Items: %v", err)
		return
	}

	if len(stressDocs.Items) > 0 {
		fmt.Printf(
			"There were %d DailyStress Items found for date range: %v - %v\n",
			len(stressDocs.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First DailyStresses ID: %s\n",
			stressDocs.Items[0].ID,
		)

		singleStressDoc, err := client.GetStress(stressDocs.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single stress item: %v", err)
			return
		}

		fmt.Printf("Single DailyStress High: %d\n", singleStressDoc.StressHigh)

	} else {
		fmt.Printf(
			"No DailyStress Items were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
