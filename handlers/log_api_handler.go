package handler

import (
	"astra/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func getNumLines(query_values url.Values) int {
	log.Println(query_values.Get("num_lines"))
	num_lines, err := strconv.Atoi(query_values.Get("num_lines"))

	if err != nil {
		log.Println(err.Error())
		return 0
	}

	if num_lines < 0 {
		return 0
	} else {
		return num_lines
	}
}

func FetchLogsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		TotalLines int
		Logs       string
		ErrorMsg   string
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

	log.Printf("NumLines: %v, FileName: %v, Filter: %v, SecondaryServer: %v", num_lines, file_name, filter, secondary_server)
	logs, err := utils.FetchLogsFromServer(num_lines, file_name, filter, secondary_server)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.ErrorMsg = err.Error()
	} else {
		resp.TotalLines = num_lines
		resp.Logs = logs
	}

	json.NewEncoder(w).Encode(&resp)
	return
}
