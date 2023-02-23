package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Abysm0xC/abysmal-api/internal/env"
	"github.com/Abysm0xC/abysmal-api/pkg/contact"
	"github.com/Abysm0xC/abysmal-api/pkg/github"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	app := fiber.New()

	app.Post("/contact", ContactHandler)
	app.Get("/healthcheck", HealthCheckHandler) // cspell: disable-line

	github := app.Group("/github")
	github.Get("/latestcommit", GithubLatestCommit) // cspell: disable-line

	return app
}

func ContactHandler(c *fiber.Ctx) error {
	// request content -> struct
	content := new(struct {
		URL     string `json:"originURL"`
		Name    string `json:"inputName"`
		Email   string `json:"inputEmail"`
		Message string `json:"inputMessage"`
	})
	err := c.BodyParser(content)
	if err != nil {
		c.Response().SetStatusCode(fiber.StatusBadRequest)
		return c.JSON(contact.Response{Success: false})
	}

	// embed
	embed := contact.Embed{
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
	body := contact.Body{
		Username:  "Contact Form",
		AvatarURL: "https://abysmal.eu.org/avatar.png",
		Embeds:    []contact.Embed{embed},
	}

	// body (struct) -> json
	jsonBody, errJson := json.Marshal(body)
	if errJson != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(contact.Response{Success: false})
	}

	// make post request
	_, errReq := http.Post(env.ContactWebhook, "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if errReq != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(contact.Response{Success: false})
	}

	return c.JSON(contact.Response{Success: true})
}

func HealthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(struct {
		Status string `json:"status"`
	}{Status: "ok"})
}

func GithubLatestCommit(c *fiber.Ctx) error {
	// request content -> struct
	content := new(struct {
		Repo   string `json:"repo"`
		User   string `json:"user,omitempty"`
		Branch string `json:"branch,omitempty"`
	})
	err := c.BodyParser(content)
	if err != nil {
		c.Response().SetStatusCode(fiber.StatusBadRequest)
		return c.JSON(github.Response{Success: false})
	}

	if len(content.User) == 0 {
		content.User = "Abysm0xC"
	}
	if len(content.Branch) == 0 {
		content.Branch = "main"
	}

	// get latest commit
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", content.User, content.Repo, content.Branch))
	if err != nil || res.StatusCode != 200 {
		c.Response().SetStatusCode(fiber.StatusBadGateway)
		return c.JSON(github.Response{Success: false})
	}
	defer res.Body.Close()

	// response body -> struct
	body, _ := io.ReadAll(res.Body)
	var resJson github.GithubResponse
	json.Unmarshal(body, &resJson)

	// take the first 7 characters
	hash := resJson.Object.SHA[0:7]

	return c.JSON(github.Response{Success: true, Hash: hash})
}
