package dto

type Bot struct {
	Token   string `json:"token" binding:"required"`
	BotName string `json:"bot_name"`
}
