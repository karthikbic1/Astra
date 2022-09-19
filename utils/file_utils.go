package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	BasePath = "/var/log"
)

func FetchLogsFromServer(num_lines int, file_name string, filter string, secondary_server string) (string, error) {

	logs := ""
	current_line := 0

	logfile, err := os.Open(fmt.Sprintf("/%v/%v", BasePath, file_name))
	defer logfile.Close()

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(logfile)

	for scanner.Scan() {
		new_line := scanner.Text()
		append := true

		if filter != "" {
			append = strings.Contains(new_line, filter)
		}

		if append {
			logs = new_line + "\n" + logs
			current_line = current_line + 1
		}

		if num_lines > 0 && num_lines == current_line {
			break
		}

	}

	return logs, nil
}
