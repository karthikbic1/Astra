package handler

import (
	"astra/core"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func MockFetchLogsFromServerSuccess(base_path string, file_name string, num_lines int, filter string, secondary_server string) (string, error) {
	return "mocked_log_content", nil
}

func MockFetchLogsFromServerFailure(base_path string, file_name string, num_lines int, filter string, secondary_server string) (string, error) {
	return "", errors.New("mock error")
}

func Test_getNumLines(t *testing.T) {
	type args struct {
		query_values url.Values
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestgetNumLines",
			args: args{query_values: url.Values{"num_lines": []string{"25"}}},
			want: 25,
		},

		{
			name: "TestgetMinNumLines",
			args: args{query_values: url.Values{"num_lines": []string{"0"}}},
			want: 1,
		},

		{
			name: "TestgetMaxNumLines",
			args: args{query_values: url.Values{"num_lines": []string{"100000000"}}},
			want: 10000,
		},
		{
			name: "TestgetDefaultNumLines",
			args: args{query_values: url.Values{}},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumLines(tt.args.query_values); got != tt.want {
				t.Errorf("getNumLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchLogsHandler(t *testing.T) {
	type args struct {
		file_name        string
		filter           string
		num_line         string
		secondary_server string
	}
	tests := []struct {
		name           string
		args           args
		expected       string
		expectedStatus int
		mockFetchLogs  func(base_path string, file_name string, num_lines int, filter string, secondary_server string) (string, error)
	}{
		{
			name:           "TestFetchLogsHandlerWithFileName",
			expected:       `{"Logs":"","ErrorMsg":"file_name request parameter is requried"}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:           "TestFetchLogsHandlerWithFileDoesntExist",
			expected:       `{"Logs":"","ErrorMsg":"Primary Server:open /var/log/dontexits: no such file or directory, Secondary Server:Not fetching from secondary server"}`,
			expectedStatus: http.StatusNotFound,
			args:           args{file_name: "dontexits"},
		},

		{
			name:           "TestFetchLogsHandlerSuccess",
			expected:       `{"Logs":"mocked_log_content","ErrorMsg":""}`,
			expectedStatus: http.StatusOK,
			args:           args{file_name: "exists"},
			mockFetchLogs:  MockFetchLogsFromServerSuccess,
		},

		{
			name:           "TestFetchLogsHandlerFailure",
			expected:       `{"Logs":"","ErrorMsg":"mock error"}`,
			expectedStatus: http.StatusInternalServerError,
			args:           args{file_name: "exists"},
			mockFetchLogs:  MockFetchLogsFromServerFailure,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/fetchlogs", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			if tt.args.file_name != "" {
				q.Add("file_name", tt.args.file_name)
			}

			if tt.args.filter != "" {
				q.Add("filter", tt.args.filter)
			}

			if tt.args.num_line != "" {
				q.Add("num_lines", tt.args.num_line)
			}

			if tt.args.secondary_server != "" {
				q.Add("secondary_server", tt.args.secondary_server)
			}

			if tt.mockFetchLogs != nil {
				// Mock FetchLogsFromServer for successful request
				orignial := core.FetchLogsFromServer
				core.FetchLogsFromServerFunc = tt.mockFetchLogs
				defer func() { core.FetchLogsFromServerFunc = orignial }()
			}

			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(FetchLogsHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expected)
			}
		})
	}
}
