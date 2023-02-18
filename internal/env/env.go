package env

import "os"

var (
	Port           string = os.Getenv("PORT")
	ContactWebhook string = os.Getenv("CONTACT_WEBHOOK")
)
