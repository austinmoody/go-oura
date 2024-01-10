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

	sleepDocs, err := client.GetSleepDocuments(threeDaysAgo, oneDaysAgo)
	if err != nil {
		fmt.Printf("Error getting Sleep Documents: %v", err)
		return
	}

	if len(sleepDocs.Documents) > 0 {
		fmt.Printf(
			"There were %d Sleep Documents found for date range: %v - %v\n",
			len(sleepDocs.Documents),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First SleepDocuments ID: %s\n",
			sleepDocs.Documents[0].ID,
		)

		singleSleepDoc, err := client.GetSleepDocument(sleepDocs.Documents[0].ID)
		if err != nil {
			fmt.Printf("Error getting single sleep document: %v", err)
			return
		}

		// Match the field 'Score' in SleepDocument struct or replace it with the correct field
		fmt.Printf("Single Sleep Document Score: %d\n", singleSleepDoc.Score)

	} else {
		fmt.Printf(
			"No Sleep Documents were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
