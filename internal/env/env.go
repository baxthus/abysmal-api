package env

import "os"

var (
	Port           string = os.Getenv("PORT")
	ContactWebhook string = os.Getenv("CONTACT_WEBHOOK")
	AvatarURL      string = os.Getenv("AVATAR_URL")
	LoggerWebhook  string = os.Getenv("LOGGER_WEBHOOK")
)
