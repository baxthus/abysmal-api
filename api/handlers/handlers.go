package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Abysm0xC/abysmal-api/internal/env"
	"github.com/Abysm0xC/abysmal-api/pkg/contact"
	"github.com/Abysm0xC/abysmal-api/pkg/github"
	ipLogger "github.com/Abysm0xC/abysmal-api/pkg/ipLogger"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	app := fiber.New()

	app.Post("/contact", ContactHandler)
	app.Get("/healthcheck", HealthCheckHandler) // cspell: disable-line
	app.Get("/iplogger", IpLoggerHandler)

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

func IpLoggerHandler(c *fiber.Ctx) error {
	// get ip
	ip, _, err := net.SplitHostPort(c.Context().RemoteAddr().String())
	if err != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(ipLogger.Response{Success: false})
	}

	// get ip info
	res, err := http.Get(fmt.Sprintf("https://ipwho.is/%s", ip))
	if err != nil || res.StatusCode != 200 {
		c.Response().SetStatusCode(fiber.StatusBadGateway)
		return c.JSON(ipLogger.Response{Success: false})
	}

	// res (json) -> struct
	body, _ := io.ReadAll(res.Body)
	var data ipLogger.Whois
	json.Unmarshal(body, &data)

	fmt.Println(data.IP)

	if !data.Success {
		c.Response().SetStatusCode(fiber.StatusBadGateway)
		return c.JSON(ipLogger.Response{Success: false})
	}

	mainField := fmt.Sprintf("**IP:** %s\n**Type:** %s", data.IP, data.Type)

	locationField := fmt.Sprintf(
		"**Continent:** %s\n**Country:** %s :flag_%s:\n**Region:** %s\n**City:** %s\n**Latitude:** %v\n**Longitude:** %v\n**Postal** %s",
		data.Continent, data.Country, strings.ToLower(data.CountryCode), data.Region, data.City, data.Latitude, data.Longitude, data.Postal,
	)

	connectionField := fmt.Sprintf(
		"**ASN:** %v\n**Organization:** %s\n**ISP:** %s",
		data.Connection.ASN, data.Connection.Org, data.Connection.ISP,
	)

	timezoneField := fmt.Sprintf(
		"**ID:** %s\n**UTC:** %s\n**Offset:** %v",
		data.Timezone.ID, data.Timezone.UTC, data.Timezone.Offset,
	)

	embed := ipLogger.Embed{
		Title:       "__Abysmal IP Logger__",
		Description: fmt.Sprintf("From: %s", c.BaseURL()),
		Thumbnail:   ipLogger.EmbedThumbnail{URL: env.AvatarURL},
		Color:       13346551,
		Fields: []ipLogger.EmbedField{
			{
				Name:  ":zap: **Main**",
				Value: mainField,
			},
			{
				// cspell: disable-next-line
				Name:  ":earth_americas: **Location**",
				Value: locationField,
			},
			{
				Name:  ":satellite: **Connection**",
				Value: connectionField,
			},
			{
				Name:  ":clock1: **Timezone**",
				Value: timezoneField,
			},
		},
	}

	button := ipLogger.ComponentsRow{
		Type: 1,
		Components: []ipLogger.Component{{
			Style:    5,
			Label:    "Open location in Google Maps",
			URL:      fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%v,%v", data.Latitude, data.Longitude),
			Disabled: false,
			Type:     2,
		}},
	}

	webhook := ipLogger.Webhook{
		Username:   "Abysmal IP Logger",
		AvatarURL:  env.AvatarURL,
		Embeds:     []ipLogger.Embed{embed},
		Components: []ipLogger.ComponentsRow{button},
	}

	// webhook (struct) -> json
	jsonWebhook, err := json.Marshal(webhook)
	if err != nil {
		c.Response().SetStatusCode(fiber.StatusInternalServerError)
		return c.JSON(ipLogger.Response{Success: false})
	}

	// make post request
	_, err = http.Post(env.LoggerWebhook, "application/json", bytes.NewBuffer([]byte(jsonWebhook)))
	if err != nil {
		c.Response().SetStatusCode(fiber.StatusBadGateway)
		return c.JSON(ipLogger.Response{Success: false})
	}

	// yes, that is what you think it is
	return c.Redirect("https://youtu.be/dQw4w9WgXcQ", fiber.StatusSeeOther)
}
