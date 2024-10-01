package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hasura/go-graphql-client"
)

const SCHEME = "https"
const SQUADCAST_HOST = "api.eu.squadcast.com"
const SQUADCAST_OAUTH_HOST = "auth.eu.squadcast.com"
const ACCESS_TOKEN_PATH = "oauth/access-token"
const GRAPHQL_PATH = "v3/graphql"
const TEAMS_PATH = "/v3/teams"
const ONCALL_PATH = "/v4/schedules/who-is-oncall"

type OnCall struct {
	RefreshToken     string
	ScheduleName     string
	TeamName         string
	TeamID           string
	SlackWebhookURL  string
	AccessToken      string
	ScheduleID       int
	OnCallShiftType  string
	OnCallResponders []string
}

type AuthenticatedClient struct {
	httpClient *http.Client
	token      string
}

// Do adds the Bearer token to each request
func (ac *AuthenticatedClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+ac.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "go-graphql-client")
	return ac.httpClient.Do(req)
}

func NewOnCall(params Params) *OnCall {
	return &OnCall{
		RefreshToken:    params.RefreshToken,
		ScheduleName:    params.ScheduleName,
		SlackWebhookURL: params.SlackWebhookURL,
		TeamName:        params.TeamName,
	}
}

func (oc *OnCall) GetAccessToken() *OnCall {
	var resp AccessTokenResponse
	var sqerr SquadcastErrorResponse

	err := NewRequest().
		Get((&url.URL{
			Scheme: SCHEME,
			Host:   SQUADCAST_OAUTH_HOST,
			Path:   ACCESS_TOKEN_PATH,
		}).String()).
		SetHeader("X-Refresh-Token", oc.RefreshToken).
		With(&resp).
		WithFail(&sqerr).
		Do()

	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	oc.AccessToken = resp.Data.AccessToken
	return oc
}

func (oc *OnCall) GetTeamID() *OnCall {
	var resp TeamsResponse
	var sqerr SquadcastErrorResponse

	err := NewRequest().
		Get((&url.URL{
			Scheme: SCHEME,
			Host:   SQUADCAST_HOST,
			Path:   TEAMS_PATH,
		}).String()).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", oc.AccessToken)).
		With(&resp).
		WithFail(&sqerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	for _, team := range resp.Data {
		if team.Name == oc.TeamName {
			oc.TeamID = team.ID
			break
		}
	}

	if oc.TeamID == "" {
		log.Fatalf("Team of name: %s doesn't exist", oc.TeamName)
	}
	return oc
}

// Description: This method is used to execute GraphQL queries by making use of the graphql package at: https://github.com/hasura/go-graphql-client
// The method accepts the GraphQL query, variables, access token and response struct as input and returns an error if the query execution fails
// Write a function executeGraphQLQuery() by following the above instructions
func executeGraphQLQuery(query interface{}, variables map[string]any, accessToken string) error {

	// AuthenticatedClient is a custom HTTP client that adds the Bearer token to the request
	// Create a new authenticated HTTP client
	authClient := &AuthenticatedClient{
		httpClient: http.DefaultClient,
		token:      accessToken,
	}

	client := graphql.NewClient(fmt.Sprintf("%s://%s/%s", SCHEME, SQUADCAST_HOST, GRAPHQL_PATH), authClient)

	// Execute the query with the Bearer token
	err := client.Query(context.Background(), query, variables)
	return err
}

//Description: As part of this method, we are fetching the schedule ID based on the team ID and schedule name
//Input: Team ID and Schedule Name
//Output: Schedule ID
//Return: OnCall object
//Steps:
//1. Execute the GraphQL query at the url: https://api.eu.squadcast.com/v3/graphql to fetch the schedule ID based on the team ID and schedule name
//2. The GraphQL query which will fetch the correct schedule id is as follows:
//{

// schedules(
//
//	filters: {
//		teamID: "66827e58f344431016e4c2de"
//		scheduleName: "team-911-pg-on-call-primary-non-office-hours"
//	}
//
//	) {
//		name
//		ID
//	}
//
// }
// 3. The response from the the above GraphQL query is as follows:
//
//	//{
//	    "data": {
//	        "schedules": [
//	            {
//	                "name": "team-911-pg-on-call-primary-non-office-hours",
//	                "ID": 1936
//	            }
//	        ]
//	    }
//	}
func (oc *OnCall) GetScheduleID() *OnCall {
	// ScheduleQuery defines the structure of the GraphQL query response
	type ScheduleQuery struct {
		Schedules []struct {
			Name string `graphql:"name"`
			ID   int    `graphql:"ID"`
		} `graphql:"schedules(filters: {teamID:$teamID, scheduleName:$scheduleName})"`
	}

	var query ScheduleQuery

	variables := map[string]any{
		"teamID":       graphql.String(oc.TeamID),
		"scheduleName": graphql.String(oc.ScheduleName),
	}

	err := executeGraphQLQuery(&query, variables, oc.AccessToken)
	if err != nil {
		log.Fatalf("There is an error fetching the schedule id for the schedule name:%s", oc.ScheduleName)
		log.Fatalf("Error:%v", err)
	}

	for _, sch := range query.Schedules {
		if sch.Name == oc.ScheduleName {
			oc.ScheduleID = sch.ID
			break
		}
	}

	if oc.ScheduleID == 0 {
		log.Fatalf("Schedule of name: %s doesn't exist", oc.ScheduleName)
	}
	return oc
}

func (oc *OnCall) GetOnCallPeople() *OnCall {
	var resp OnCallApiResponse
	var sqerr SquadcastErrorResponse

	baseURL := &url.URL{
		Scheme: SCHEME,
		Host:   SQUADCAST_HOST,
		Path:   ONCALL_PATH,
	}
	queryParams := url.Values{}
	queryParams.Add("teamId", oc.TeamID)
	queryParams.Add("scheduleID", strconv.Itoa(oc.ScheduleID)) // Correctly convert scheduleID to string
	baseURL.RawQuery = queryParams.Encode()

	err := NewRequest().
		Get(baseURL.String()).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", oc.AccessToken)).
		With(&resp).
		WithFail(&sqerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	for _, data := range resp.Data {
		for _, onCallResponder := range data.Oncall {
			oc.OnCallResponders = append(oc.OnCallResponders, fmt.Sprintf("%s %s", onCallResponder.FirstName, onCallResponder.LastName))
		}
	}
	return oc
}

func (oc *OnCall) NotifySlack() {
	var resp string
	var slkerr string

	fields := make([]SlackWebhookAttachmentFields, 0)
	for _, user := range oc.OnCallResponders {
		fields = append(fields, SlackWebhookAttachmentFields{
			Title: user,
		})
	}

	body := SlackWebhookRequest{
		Attachments: []SlackWebhookAttachment{
			{
				Fallback: fmt.Sprintf("On-Call Update for schedule: %s", oc.ScheduleName),
				Pretext:  fmt.Sprintf("People on-call for schedule: %s", oc.ScheduleName),
				Color:    "#00FF00",
				Fields:   fields,
			},
		},
	}

	err := NewRequest().
		Post(oc.SlackWebhookURL).
		Data(body).
		With(&resp).
		WithFail(&slkerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, slkerr)
	}
}
