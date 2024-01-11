package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	startDateTIme := time.Now().Add(-25 * time.Hour)
	endDateTime := time.Now().Add(-24 * time.Hour)

	heartrates, err := client.GetHeartRates(startDateTIme, endDateTime)
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
			"First HeartRate BPM is %d Source is %s at %v",
			heartrates.Items[0].Bpm,
			heartrates.Items[0].Source,
			heartrates.Items[0].Timestamp.Format("02-Jan-2006 15:04:05"),
		)
	}
}
