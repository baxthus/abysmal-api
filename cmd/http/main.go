package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	port           string = os.Getenv("PORT")
	contactWebhook string = os.Getenv("CONTACT_WEBHOOK")
)

type contactResponse struct {
	Success bool `json:"success"`
}

type contactEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int32  `json:"color"`
	Fields      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
}

type contactBody struct {
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Embeds    []contactEmbed `json:"embeds"`
}

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
		return c.JSON(contactResponse{Success: false})
	}

	// embed
	embed := contactEmbed{
		Title:       "Contact Form",
		Description: fmt.Sprintf("Form %s at <t:%s:f>", content.URL, strconv.Itoa(int(time.Now().Unix()))),
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
				Value: fmt.Sprintf("`%s", content.Email),
			},
			{
				Name:  ":page_facing_up: **Message**",
				Value: fmt.Sprintf("```%s```", content.Message),
			},
		},
	}

	// request body
	body := contactBody{
		Username:  "Contact Form",
		AvatarURL: "https://abysmal.eu.org/avatar.png",
		Embeds:    []contactEmbed{embed},
	}

	// body (struct) -> json
	jsonBody, errJson := json.Marshal(body)
	if errJson != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(contactResponse{Success: false})
	}

	// make post request
	_, errReq := http.Post(contactWebhook, "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if errReq != nil {
		c.Response().Header.SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(contactResponse{Success: false})
	}

	return c.JSON(contactResponse{Success: true})
}

func main() {
	app := fiber.New()

	app.Static("/", "/public")
	app.Post("/contact", contactHandler)

	app.Listen("0.0.0.0:" + port)
}
