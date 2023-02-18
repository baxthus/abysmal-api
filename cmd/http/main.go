package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Abysm0xC/abysmal-api/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	port           string = os.Getenv("PORT")
	contactWebhook string = os.Getenv("CONTACT_WEBHOOK")
)

func contactHandler(c *fiber.Ctx) error {
	// request content -> struct
	content := new(struct {
		URL     string `json:"originURL"`
		Name    string `json:"inputName"`
		Email   string `json:"inputEmail"`
		Message string `json:"inputMessage"`
	})
	err := c.BodyParser(content)
	if err != nil {
		c.Response().Header.SetStatusCode(fiber.StatusBadRequest)
		return c.JSON(types.ContactResponse{Success: false})
	}

	// embed
	embed := types.ContactEmbed{
		Title:       "Contact Form",
		Description: fmt.Sprintf("Form %s at <t:%v:f>", content.URL, time.Now().Unix()),
		Color:       13346551,
		Fields: []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			{
				Name:  ":bust_in_silhouette: **Name**",
				Value: fmt.Sprintf("`%s`", content.Name),
			},
			{
				Name:  ":envelope: **Email**",
				Value: fmt.Sprintf("`%s`", content.Email),
			},
			{
				Name:  ":page_facing_up: **Message**",
				Value: fmt.Sprintf("```%s```", content.Message),
			},
		},
	}

	// request body
	body := types.ContactBody{
		Username:  "Contact Form",
		AvatarURL: "https://abysmal.eu.org/avatar.png",
		Embeds:    []types.ContactEmbed{embed},
	}

	// body (struct) -> json
	jsonBody, errJson := json.Marshal(body)
	if errJson != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(types.ContactResponse{Success: false})
	}

	// make post request
	_, errReq := http.Post(contactWebhook, "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if errReq != nil {
		c.Response().Header.SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(types.ContactResponse{Success: false})
	}

	return c.JSON(types.ContactResponse{Success: true})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/", "/public")
	app.Post("/contact", contactHandler)

	app.Get("/healthcheck", func(c *fiber.Ctx) error { // cspell: disable-line
		return c.JSON(struct {
			Status string `json:"status"`
		}{Status: "ok"})
	})

	app.Listen("0.0.0.0:" + port)
}
