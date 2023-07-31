package main

import (
	"fmt"
	"github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"overseerr-tautulli-data-exporter/config"
	"path/filepath"
)

func main() {
	// Init logging
	logFilePath := filepath.Clean(fmt.Sprintf("%s/%s", os.Getenv("LOG_DIR"), "overseerr-tautulli-data-exporter.log"))
	logFile, errLog := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0666,
	)
	if errLog != nil {
		log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Warn().Stack().Err(errLog).Str("file", logFilePath).Msg("Error opening log file, will write on console only")
	} else {
		log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile))
		log.Info().Str("file", logFilePath).Msg(fmt.Sprintf("Logging to file %s", logFile.Name()))
	}

	log.Info().Str("version", os.Getenv("DOCKER_TAG")).Msg("overseerr-tautulli-data-exporter starting")

	// Reading config
	config.LoadConfig()

	// Defining gin server
	router := gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Error().Stack().Err(err).Msg("router.SetTrustedProxies returned an error")
		os.Exit(1)
	}
	router.Use(ginzerolog.Logger("gin"))

	// TODO init routes here
	router.GET("/check", check)
	router.GET("/version", version)

	// Defining gin server
	routerErr := router.Run(":8090")
	if routerErr != nil {
		log.Error().Stack().Err(routerErr).Msg("router.Run returned an error")
		os.Exit(1)
	}
}

func check(context *gin.Context) {
	context.String(http.StatusOK, "OK")
}

func version(context *gin.Context) {
	version := os.Getenv("DOCKER_TAG")
	if len(version) < 1 {
		version = "DOCKER_TAG not defined"
	}
	context.String(http.StatusOK, fmt.Sprintf("Version: %s", version))
}
