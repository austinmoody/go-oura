package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	threeDaysAgo := time.Now().AddDate(-1, 0, 0)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	sessionDocs, err := client.GetSessions(threeDaysAgo, oneDaysAgo, nil)

	if err != nil {
		fmt.Printf("Error getting Session Items: %v", err)
		return
	}

	if len(sessionDocs.Items) > 0 {
		fmt.Printf(
			"There were %d Session Items found for date range: %v - %v\n",
			len(sessionDocs.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Session ID: %s\n",
			sessionDocs.Items[0].ID,
		)

		singleSessionDoc, err := client.GetSession(sessionDocs.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single session item: %v", err)
			return
		}

		fmt.Printf("Single Session Type: %s\n", singleSessionDoc.Type)

	} else {
		fmt.Printf(
			"No Session Items were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
