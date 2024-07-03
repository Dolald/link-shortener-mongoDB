package main

import (
	"context"
	"shortener/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
