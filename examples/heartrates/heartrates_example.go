package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	startDateTIme := time.Now().Add(-168 * time.Hour) // crank up to see next_token
	endDateTime := time.Now().Add(-24 * time.Hour)

	heartrates, err := client.GetHeartRates(startDateTIme, endDateTime, nil)
	if err != nil {
		fmt.Printf("Error getting heartrates: %v", err)
		return
	}

	if len(heartrates.Items) > 0 {
		fmt.Printf(
			"There were %d HearRates found for date range: %v - %v\n",
			len(heartrates.Items),
			startDateTIme.Format("02-Jan-2006 15:04:05"),
			endDateTime.Format("02-Jan-2006 15:04:06"),
		)

		fmt.Printf(
			"First HeartRate BPM is %d Source is %s at %v\n",
			heartrates.Items[0].Bpm,
			heartrates.Items[0].Source,
			heartrates.Items[0].Timestamp.Format("02-Jan-2006 15:04:05"),
		)

		if heartrates.NextToken != "" {
			fmt.Printf("\tThere was a next_token found of %s\n", heartrates.NextToken)
			heartrates, err = client.GetHeartRates(startDateTIme, endDateTime, &heartrates.NextToken)
			if err != nil {
				fmt.Printf("Error getting heartrates: %v", err)
				return
			}

			if len(heartrates.Items) > 0 {
				fmt.Printf("\tPulled %d more heartrates with next token\n", len(heartrates.Items))
			} else {
				fmt.Printf("\tNo more heartrates found with next token\n")
			}

		} else {
			fmt.Printf("\tThere was NO next_token in the response\n")
		}
	}
}
