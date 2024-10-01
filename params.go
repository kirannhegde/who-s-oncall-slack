package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Params struct contains all external parameters
// for getting on-call people
type Params struct {
	RefreshToken    string `json:"REFRESH_TOKEN"`
	ScheduleName    string `json:"SCHEDULE_NAME"`
	SlackWebhookURL string `json:"SLACK_WEBHOOK_URL"`
	TeamName        string `json:"SQUADCAST_TEAM_NAME"`
}

var params = Params{}

func loadFromJSON(files ...string) {
	for _, file := range files {
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] unable to load from config file: %s\n", err)
			continue
		}

		err = json.Unmarshal(bs, &params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] unable to load from config file: %s\n", err)
			return
		}
	}
}
