package main

import (
	"fmt"
	"gas-rest-api/internal/config"
)

func main() {
	// Reading config file.
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger (slog)

	// TODO: init storage

	// TODO: init router

	// TODO: run server

}
