package main

import (
	"fmt"
	"github.com/austinmoody/go-oura"
	"os"
	"time"
)

func main() {
	client := go_oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	/*
		readiness, _ := client.GetReadinessDocument("29a809a2-778c-4742-b945-e01876b8f32a")
		fmt.Println(readiness.Id)
	*/

	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	oneDaysAgo := time.Now().Add(-24 * time.Hour)

	readiness, err := client.GetReadinessDocuments(twoDaysAgo, oneDaysAgo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(readiness.Documents))
}
