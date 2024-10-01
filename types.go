package main

type SquadcastErrorResponse struct {
	Meta SquadcastErrorDetails `json:"meta"`
}

type SquadcastErrorDetails struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type AccessTokenResponse struct {
	Data AccessTokenDetails `json:"data"`
}

type AccessTokenDetails struct {
	AccessToken string `json:"access_token"`
}

type SchedulesResponse struct {
	Data []SchedulesDetails `json:"data"`
}

type SchedulesDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamsResponse struct {
	Data []TeamsDetails `json:"data"`
}

type TeamsDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Structs to represent the JSON structure of the response
type Contact struct {
	DialCode    string `json:"dial_code"`
	PhoneNumber string `json:"phone_number"`
}

type OnCallPerson struct {
	ID                 string  `json:"id"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	UsernameForDisplay string  `json:"username_for_display"`
	Email              string  `json:"email"`
	Contact            Contact `json:"contact"`
}

type Rotation struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
}

type Schedule struct {
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	DeletedAt string     `json:"deleted_at"`
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Rotations []Rotation `json:"rotations"`
}

type Data struct {
	Schedule Schedule       `json:"schedule"`
	Oncall   []OnCallPerson `json:"oncall"`
}

type OnCallApiResponse struct {
	Data []Data `json:"data"`
}

type SlackWebhookRequest struct {
	Text        string                   `json:"text,omitempty"`
	Channel     string                   `json:"channel,omitempty"`
	Username    string                   `json:"username,omitempty"`
	IconURL     string                   `json:"icon_url,omitempty"`
	IconEmoji   string                   `json:"icon_emoji,omitempty"`
	Attachments []SlackWebhookAttachment `json:"attachments,omitempty"`
}

type SlackWebhookAttachment struct {
	Fallback string                         `json:"fallback,omitempty"`
	Pretext  string                         `json:"pretext,omitempty"`
	Color    string                         `json:"color,omitempty"`
	Fields   []SlackWebhookAttachmentFields `json:"fields,omitempty"`
}

type SlackWebhookAttachmentFields struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}
