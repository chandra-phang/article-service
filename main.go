package main

import (
	"flag"

	"article-service/app"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.yml", "absolute path to the configuration file")
	flag.Parse()

	application := app.NewApplication()
	application.InitApplication(configFilePath)
}
