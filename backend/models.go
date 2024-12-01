package main

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

type AccessToken struct {
	Token string `json:"access_token"`
}

type TokenParams struct {
	DiscordID string `json:"discord_id"`
	ClientID  string `json:"client_id"`
}
