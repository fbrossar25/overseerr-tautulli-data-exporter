package main

import (
	"fmt"
	"github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var MaxLogLevel zerolog.Level
var AppLogger zerolog.Logger

type LogLevelHook struct{}

func (h LogLevelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level < MaxLogLevel {
		e.Discard()
	}
}

func main() {
	// Init logging
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(strings.TrimSpace(logLevel)) {
	case "ALL", "TRACE":
		MaxLogLevel = zerolog.TraceLevel
		MaxLogLevel = zerolog.TraceLevel
	case "DEBUG":
		MaxLogLevel = zerolog.DebugLevel
	case "INFO":
		MaxLogLevel = zerolog.InfoLevel
	case "WARN":
		MaxLogLevel = zerolog.WarnLevel
	case "ERROR":
		MaxLogLevel = zerolog.ErrorLevel
	case "DISABLED":
		MaxLogLevel = zerolog.Disabled
	default:
		MaxLogLevel = zerolog.InfoLevel
	}
	logFilePath := filepath.Clean(fmt.Sprintf("%s/%s", os.Getenv("LOG_DIR"), "overseerr-tautulli-data-exporter.log"))
	logFile, errLog := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0666,
	)
	if errLog != nil {
		AppLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Hook(LogLevelHook{})
		AppLogger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		AppLogger.Warn().Stack().Err(errLog).Str("file", logFilePath).Msg("Error opening log file, will write on console only")
	} else {
		AppLogger = zerolog.New(io.MultiWriter(logFile, zerolog.ConsoleWriter{Out: os.Stderr})).With().Timestamp().Logger().Hook(LogLevelHook{})
		AppLogger.Info().Str("file", logFilePath).Msg(fmt.Sprintf("Logging to file %s", logFile.Name()))
	}

	AppLogger.Info().Str("version", os.Getenv("DOCKER_TAG")).Str("mexLogLevel", MaxLogLevel.String()).Msg("overseerr-tautulli-data-exporter starting")

	// Reading config
	LoadConfig()

	// Defining gin server
	router := gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		AppLogger.Error().Stack().Err(err).Msg("router.SetTrustedProxies returned an error")
		os.Exit(1)
	}
	router.Use(ginzerolog.Logger("gin"))

	// TODO init routes here
	router.GET("/check", check)
	router.GET("/version", version)
	InitTautulli()

	// Defining gin server
	routerErr := router.Run(":8090")
	if routerErr != nil {
		AppLogger.Error().Stack().Err(routerErr).Msg("router.Run returned an error")
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
