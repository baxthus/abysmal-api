package main

import (
	"github.com/Abysm0xC/abysmal-api/api/router"
	"github.com/Abysm0xC/abysmal-api/internal/env"
)

func main() {
	app := router.Initialize()

	app.Listen("0.0.0.0:" + env.Port)
}
