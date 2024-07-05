package main

import (
	"context"
	"shortener/internal/app"
)

func main() {
	app.Run(context.Background())
}
