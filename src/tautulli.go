package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TautulliResponse[T interface{}] struct {
	Response struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		Data    T      `json:"data"`
	} `json:"response"`
}

func InitTautulli() {
	AppLogger.Info().Str("url", Config.Tautulli.Url).Str("apikey", fmt.Sprintf("%s***", Config.Tautulli.ApiKey[0:3])).Msg("Tautulli exporter starting")
	TautulliCheck()
}

func getRedactedQuery(response *http.Response) url.Values {
	query := response.Request.URL.Query()
	if query.Has("apikey") {
		query.Set("apikey", fmt.Sprintf("%s***", query.Get("apikey")[0:3]))
	}
	return query
}

func GETTautulli(cmd string, target interface{}) error {
	response, err := http.Get(
		fmt.Sprintf("%s/api/v2?cmd=%s&apikey=%s", Config.Tautulli.Url, cmd, Config.Tautulli.ApiKey))
	if err != nil {
		AppLogger.Error().
			Int("statusCode", response.StatusCode).
			Str("requestedUrl", response.Request.URL.Path).
			Interface("query", getRedactedQuery(response)).
			Msg("Error calling Tautulli (GET)")
		return err
	} else {
		AppLogger.Debug().
			Int("statusCode", response.StatusCode).
			Str("requestedUrl", response.Request.URL.Path).
			Interface("query", getRedactedQuery(response)).
			Msg("Success calling Tautulli (GET)")
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

func TautulliCheck() bool {
	status := new(TautulliResponse[interface{}])
	err := GETTautulli("status", status)
	if err != nil {
		AppLogger.Error().Str("url", Config.Tautulli.Url).Msg("Error checking Tautulli status")
		return false
	} else if strings.ToUpper(strings.TrimSpace(status.Response.Message)) != "OK" {
		AppLogger.Warn().Str("url", Config.Tautulli.Url).Interface("response", status.Response).Msg("Tautulli not available")
		return false
	} else {
		AppLogger.Info().Interface("response", status).Str("url", Config.Tautulli.Url).Msg("Tautulli available")
		return true
	}
}
