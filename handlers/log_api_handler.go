package handler

import (
	"astra/core"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	BasePath       = "/var/log"
	MinLogLines    = 1
	MaxLogLines    = 10000
	DefaultLogLine = 10
)

func getNumLines(query_values url.Values) int {
	// This method ensures that proper num_lines is passed to fetch the logs.
	// Adds the guard rails for the parameter.
	num_lines, err := strconv.Atoi(query_values.Get("num_lines"))

	// Any error parsing the num_lines parameter, then use the MinLogLines
	if err != nil {
		return DefaultLogLine
	}

	// if user provides num_lines less than MinLogLines, then ensure we return atleast MinLogLines
	if num_lines <= 0 {
		return MinLogLines
		// if user provides num_lines greater MaxLogLines, ensure that we dont cross the MaxLogLines.
	} else if num_lines > MaxLogLines {
		return MaxLogLines
	} else {
		// if num_lines is within bounds, then return that
		return num_lines
	}
}

func FetchLogsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Logs     string
		ErrorMsg string
	}

	resp := &response{}

	// Read the params from the Query URL
	query_values := r.URL.Query()
	num_lines := getNumLines(query_values)
	file_name := query_values.Get("file_name")
	filter := query_values.Get("filter")
	secondary_server := query_values.Get("secondary_server")

	// if user has not provided file_name, send a 400 badrequest response.
	if file_name == "" {
		resp.ErrorMsg = "file_name request parameter is requried"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	// Fetch the logs in latest in fashion with user provided conditions.
	logs, err := core.FetchLogsFromServer(BasePath, file_name, num_lines, filter, secondary_server)

	// Error processing the file return proper error to the user.
	if err != nil {

		resp.ErrorMsg = err.Error()
		if strings.Contains(resp.ErrorMsg, "no such file or directory") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		log.Println(resp.ErrorMsg)

	} else {
		resp.Logs = logs
	}

	// If we reach here, it means no error and return proper response with 200 Ok.
	json.NewEncoder(w).Encode(&resp)
	return
}
