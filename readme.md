# Who's Oncall Slack

The purpose of the script is to get the list of people who are on-call for a given Squadcast Schedule across all the shifts, at the time of execution of the script and send it to a slack channel.

## Script Inputs

The script takes a configuration JSON file as its command-line argument which looks something like the following:

```json
{
  "REFRESH_TOKEN": "Squadcast Refresh Token which will be used to be make requests to Public APIs",
  "SCHEDULE_NAME": "Name of the schedule for which we are getting the on-call people",
  "SLACK_WEBHOOK_URL": "Slack incoming webhook URL for the channel to which we want to send the on-call notification",
  "SQUADCAST_TEAM_NAME": "Squadcast team where the on-call schedule is present."
}
```

Update the given `params.json` with suitable values.

## Building and Running

```bash
go build && ./oncall params.json
```
