package main

import (
	"context"
	"shortener/pkg/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
