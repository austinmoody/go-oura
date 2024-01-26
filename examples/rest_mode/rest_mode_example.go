package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	rightNow := time.Now().Add(-24 * time.Hour)

	restModes, err := client.GetRestModes(threeMonthsAgo, rightNow, nil)
	if err != nil {
		fmt.Printf("Error getting Rest Modes: %v", err)
		return
	}

	if len(restModes.Items) > 0 {
		fmt.Printf(
			"There were %d Rest Modes found for date range: %v - %v\n",
			len(restModes.Items),
			threeMonthsAgo.Format("02-Jan-2006"),
			rightNow.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Rest Mode ID: %s\n",
			restModes.Items[0].ID,
		)

		singleRestMode, err := client.GetRestMode(restModes.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single rest mode item: %v", err)
			return
		}

		fmt.Printf("Single Rest Mode: %v\n", singleRestMode)

	} else {
		fmt.Printf(
			"No Rest Modes were found for the date range: %v - %v",
			threeMonthsAgo.Format("02-Jan-2006"),
			rightNow.Format("02-Jan-2006"),
		)
	}
}
