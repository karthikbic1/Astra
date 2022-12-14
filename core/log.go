package core

import (
	httputils "astra/httputil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var FetchLogsFromServerFunc = FetchLogsFromServer

func FetchLogsFromSecondaryServer(num_lines int, file_name string, filter string, secondary_server string) ([]string, error) {
	if secondary_server == "" {
		return []string{}, errors.New("Not fetching from secondary server")
	}

	query := map[string]interface{}{
		// this appended to mimic the behaviour of file not present on primary server
		// but exists on the secondary server. Idea is file_name: xyz wont be present,
		// but file_name: xyz-secondary will be, so when the request hits the primary
		// server file wont be found and its forward to second server with `xyz-secondary`
		// which will exists and secondary server will return.
		"file_name": file_name + "-secondary",
		"num_lines": num_lines,
	}

	if filter != "" {
		query["filter"] = filter
	}

	status_code, bytes, err := httputils.Get(
		context.Background(),
		secondary_server+"/fetchlogs",
		query,
	)

	if status_code != http.StatusOK || err != nil {
		return []string{}, err
	}

	type jsonMap struct {
		ErrorMsg string
		Logs     []string
	}

	var resp jsonMap
	err = json.Unmarshal([]byte(string(bytes)), &resp)

	if err != nil {
		return []string{}, err
	}

	return resp.Logs, nil
}

func FetchLogsFromServer(base_path string, file_name string, num_lines int, filter string, secondary_server string) ([]string, error) {

	// This the core business logic of the fetch logs api
	// Read lines from the given file starting from EOF and
	// seek backeards to apply the conditions such no of
	// lines, filters in a single pass.

	logs_arr := []string{}
	new_line := ""
	line_count := 0
	next_char := make([]byte, 1)

	// open the file for reading
	logfile, err := os.Open(fmt.Sprintf("%v/%v", base_path, file_name))
	defer logfile.Close()

	if err != nil {
		// if file doesnt exists on the primary server
		// try to fetch it from secondary server if provided
		secondary_server_logs, ss_err := FetchLogsFromSecondaryServer(num_lines, file_name, filter, secondary_server)
		if ss_err != nil {
			ss_err = errors.New("Primary Server:" + err.Error() + ", Secondary Server:" + ss_err.Error())
		}
		return secondary_server_logs, ss_err
	}

	// Seek to the end of file starting with last character before EOF.
	logfile.Seek(-1, io.SeekEnd)

	// Calculate the cursor positions,
	current_position, _ := logfile.Seek(0, io.SeekCurrent)

	// Seek backwards until a new line character is encountered.
	// Seeks until the start of the file.
	for current_position >= -1 {

		// Seek to the new position
		logfile.Seek(current_position, io.SeekCurrent)
		// Read the character at the current position
		logfile.ReadAt(next_char, current_position)
		// Decrease the current cursor to point to next character on the left.
		current_position = current_position - 1

		// If a new line character is encountered then
		// its time to form the line and apply conditions.
		if string(next_char) == "\n" { // Append the character to create line. Since we are reading backeards

			append_str := true

			// check if filter is present and filter is part of newly formed line.
			if filter != "" {
				append_str = strings.Contains(new_line, filter)
			}

			// append the line to the log lines
			// Keep track of no of line
			if append_str && new_line != "" {

				logs_arr = append(logs_arr, new_line)
				line_count = line_count + 1

			}

			// Reset the new_line variable for next iteration to hold new line text.
			new_line = ""

			// if user has provided num_lines, then break out of the loop
			// when the counter becomes the user provided count.
			if num_lines > 0 && num_lines == line_count {
				break
			}
		} else {
			// Append the character to create line. Since we are reading backeards
			// append the character in the front of the line.
			new_line = string(next_char) + new_line
		}
	}

	return logs_arr, nil
}
