package main

import (
	"log"
	"os"
)

func main() {
	// init config/env, pass params json path
	if len(os.Args) < 2 {
		log.Fatal("params file not specified")
	}

	// load all params
	loadFromJSON(os.Args[1:]...)

	// execute pipeline
	NewOnCall(params).
		GetAccessToken().
		GetTeamID().
		GetScheduleID().
		GetOnCallPeople().
		NotifySlack()
}
