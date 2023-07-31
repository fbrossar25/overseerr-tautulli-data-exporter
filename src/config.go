package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type ExporterConfig struct {
	Tautulli struct {
		Url    string `yaml:"url" validate:"required,min=1"`
		ApiKey string `yaml:"apiKey" validate:"required,min=1"`
	} `yaml:"tautulli"`
	Overseerr struct {
		Url    string `yaml:"url" validate:"required,min=1"`
		ApiKey string `yaml:"apiKey" validate:"required,min=1"`
	} `yaml:"overseerr"`
}

var Config ExporterConfig

func LoadConfig() {
	configFilePath := filepath.Clean(fmt.Sprintf("%s/%s", os.Getenv("CONF_DIR"), "overseerr-tautulli-data-exporter.yml"))
	f, configFileErr := os.Open(configFilePath)
	if configFileErr != nil {
		log.Error().Stack().Err(configFileErr).Str("file", configFilePath).Msg("Error opening config file")
		panic(configFileErr)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	configErr := decoder.Decode(&Config)
	if configErr != nil {
		log.Error().Stack().Err(configErr).Str("file", configFilePath).Msg("Error reading config file")
		panic(configErr)
	}
	v := validator.New()
	validateConfigErr := v.Struct(Config)
	if validateConfigErr != nil {
		for _, e := range validateConfigErr.(validator.ValidationErrors) {
			log.Error().Stack().Str("configError", fmt.Sprint(e)).Msg("Erreur à la lecture du fichier de config")
		}
		panic(validateConfigErr)
	}
	log.Info().Str("file", configFilePath).Msg("Config file parsed")
}
