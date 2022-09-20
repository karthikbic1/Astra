package handler

import (
	"astra/core"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BasePath    = "/var/log"
	MinLogLines = 10
	MaxLogLines = 10000
)

func getNumLines(query_values url.Values) int {
	num_lines, err := strconv.Atoi(query_values.Get("num_lines"))

	if err != nil {
		log.Println(err.Error())
		return MinLogLines
	}

	if num_lines < MinLogLines {
		return MinLogLines
	} else if num_lines > MaxLogLines {
		return MaxLogLines
	} else {
		return num_lines
	}
}

func FetchLogsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Logs     string
		ErrorMsg string
	}

	resp := &response{}

	query_values := r.URL.Query()
	num_lines := getNumLines(query_values)
	file_name := query_values.Get("file_name")
	filter := query_values.Get("filter")
	secondary_server := query_values.Get("secondary_server")

	if file_name == "" {
		resp.ErrorMsg = "file_name request parameter is requried"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}

	logs, err := core.FetchLogsFromServer(BasePath, file_name, num_lines, filter, secondary_server)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		resp.ErrorMsg = err.Error()
	} else {
		resp.Logs = logs
	}

	json.NewEncoder(w).Encode(&resp)
	return
}
