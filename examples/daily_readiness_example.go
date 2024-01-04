package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {

	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	// Get Multiple Readiness Documents (if they exist) for the last couple days.
	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	readiness, err := client.GetReadinessDocuments(twoDaysAgo, oneDaysAgo)
	if err != nil {
		fmt.Printf("Error getting multiple readiness documents: %v", err)
		return
	}

	if len(readiness.Documents) > 0 {
		fmt.Printf(
			"There were %d Readiness Documents found for date range: %v - %v\n",
			len(readiness.Documents),
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Readiness Document ID: %s\n",
			readiness.Documents[0].Id,
		)

		fmt.Printf(
			"First Readiness Document Score: %d\n",
			readiness.Documents[0].Score,
		)
	} else {
		fmt.Printf(
			"No Readiness Documents were found for the date range: %v - %v",
			twoDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}

	/*
		readiness, _ := client.GetReadinessDocument("29a809a2-778c-4742-b945-e01876b8f32a")
		fmt.Println(readiness.Id)
	*/

}
