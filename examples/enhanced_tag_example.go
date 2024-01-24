package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	threeDaysAgo := time.Now().AddDate(0, -2, 0)
	oneDaysAgo := time.Now()

	tags, err := client.GetEnhancedTags(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting EnhancedTag Items: %v", err)
		return
	}

	if len(tags.Items) > 0 {
		fmt.Printf(
			"There were %d EnhancedTag Items found for date range: %v - %v\n",
			len(tags.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First EnhancedTag ID: %s\n",
			tags.Items[0].ID,
		)

		singleTag, err := client.GetEnhancedTag(tags.Items[0].ID)
		if err != nil {
			fmt.Printf("Error getting single tag item: %v", err)
			return
		}

		// Assuming 'TagTypeCode' to be useful piece of information from EnhancedTag struct
		fmt.Printf("Single EnhancedTag Type Code: %s\n", singleTag.TagTypeCode)

	} else {
		fmt.Printf(
			"No EnhancedTag Items were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
