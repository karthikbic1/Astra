package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func FetchLogsFromSecondaryServer(num_lines int, file_name string, filter string, secondary_server string) (string, error) {
	return "from secondary server..", errors.New("Not Implemented..")
}

func FetchLogsFromServer(base_path string, file_name string, num_lines int, filter string, secondary_server string) (string, error) {

	logs := ""
	new_line := ""
	line_count := 0
	next_char := make([]byte, 1)

	logfile, err := os.Open(fmt.Sprintf("/%v/%v", base_path, file_name))
	defer logfile.Close()

	if err != nil {
		secondary_server_logs, ss_err := FetchLogsFromSecondaryServer(num_lines, file_name, filter, secondary_server)
		if ss_err != nil {
			ss_err = errors.New("Primary Server:" + err.Error() + ", Secondary Server:" + ss_err.Error())
		}
		return secondary_server_logs, ss_err
	}

	logfile.Seek(-1, io.SeekEnd)

	current_position, _ := logfile.Seek(0, io.SeekCurrent)

	for current_position >= -1 {

		logfile.Seek(current_position, io.SeekCurrent)
		logfile.ReadAt(next_char, current_position)
		current_position = current_position - 1

		new_line = string(next_char) + new_line

		if string(next_char) == "\n" {
			append := true

			if filter != "" {
				append = strings.Contains(new_line, filter)
			}

			if append {
				logs = logs + new_line
				line_count = line_count + 1
			}

			new_line = ""

			if num_lines > 0 && num_lines == line_count {
				break
			}
		}
	}

	return logs, nil
}
