package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TautulliResponse[T interface{}] struct {
	Response struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		Data    T      `json:"data"`
	} `json:"response"`
}

type Epoch int64

// TautulliLibraryMediaInfo Response to command get_library_media_info
// body of request (form-data):
//
//	section_id : int, library id from plex
//	refresh : boolean, trigger refresh or not
//	length : int, how many elements to retrieve, can be 0
//
// RecordsTotal Total number of records in this library with current filters
// RecordsTotal Total number of records in this library
// Data Filtered medias data
// FilteredFileSize Size of all elements displayed in Data
// TotalFileSize Size of Library
// LastRefreshed Last time media infos refreshed (Epoch milliseconds)
type TautulliLibraryMediaInfo struct {
	RecordsFiltered  int   `json:recordsFiltered`
	RecordsTotal     int   `json:recordsTotal`
	Data             []any `json:data`
	Draw             int   `json:draw`
	FilteredFileSize int   `json:filtered_file_size`
	TotalFileSize    int   `json:total_file_size`
	LastRefreshed    Epoch `json:last_refreshed`
}

type LibrarySizeMetric struct {
	gorm.Model
	LibraryName   string
	PlexSectionId int
	Timestamp     time.Time
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

func GETTautulli(cmd string, target interface{}, params map[string]string) error {
	values := url.Values{}
	values.Add("cmd", cmd)
	values.Add("apikey", Config.Tautulli.ApiKey)
	if params != nil {
		for key, element := range params {
			values.Add(key, element)
		}
	}
	response, err := http.Get(
		fmt.Sprintf("%s/api/v2?%s", Config.Tautulli.Url, values.Encode()))
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
	err := GETTautulli("status", status, nil)
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
