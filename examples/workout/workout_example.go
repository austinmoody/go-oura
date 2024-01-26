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

	workoutDocs, err := client.GetWorkouts(threeDaysAgo, oneDaysAgo, nil)
	if err != nil {
		fmt.Printf("Error getting Workout Items: %v", err)
		return
	}

	if len(workoutDocs.Items) > 0 {
		fmt.Printf(
			"There were %d Workout Items found for date range: %v - %v\n",
			len(workoutDocs.Items),
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)

		fmt.Printf(
			"First Workout ID: %s\n",
			workoutDocs.Items[0].Id,
		)

		singleWorkoutDoc, err := client.GetWorkout(workoutDocs.Items[0].Id)
		if err != nil {
			fmt.Printf("Error getting single workout item: %v", err)
			return
		}

		fmt.Printf("Single Workout Activity: %s\n", singleWorkoutDoc.Activity)
		fmt.Printf("Single Workout Calories: %f\n", singleWorkoutDoc.Calories)
		fmt.Printf("Single Workout Distance: %f\n", singleWorkoutDoc.Distance)

	} else {
		fmt.Printf(
			"No Workout Items were found for the date range: %v - %v",
			threeDaysAgo.Format("02-Jan-2006"),
			oneDaysAgo.Format("02-Jan-2006"),
		)
	}
}
