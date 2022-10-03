package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFetchLogsFromServer(t *testing.T) {
	type args struct {
		base_path        string
		file_name        string
		num_lines        int
		filter           string
		secondary_server string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
		errMsg  string
	}{
		{
			name: "TestFetchLogsFromServer",
			args: args{
				base_path: ".",
				file_name: "testlog",
			},
			want: []string{"this is 9 line.",
				"this is 8 line ",
				"this is 7 line",
				"this is 6 line.",
				"this is 5 line 今日は",
				"this is 4 line",
				"this is 3 line.",
				"this is 2 line",
				"this is 1 line",
			},
			wantErr: false,
		},
		{
			name: "TestFetchLogsFromServerWithNumLines",
			args: args{
				base_path: ".",
				file_name: "testlog",
				num_lines: 2,
			},

			want: []string{"this is 9 line.",
				"this is 8 line ",
			},
			wantErr: false,
		},

		{
			name: "TestFetchLogsFromServerWithFilter",
			args: args{
				base_path: ".",
				file_name: "testlog",
				filter:    "5 line",
			},
			want: []string{
				"this is 5 line 今日は"},
			wantErr: false,
		},

		{
			name: "TestFetchLogsFromServerWithFileMissingAndNoSecondaryServer",
			args: args{
				base_path: ".",
				file_name: "doesntexists",
			},
			want:    []string{},
			wantErr: true,
			errMsg:  "Primary Server:open ./doesntexists: no such file or directory, Secondary Server:Not fetching from secondary server",
		},

		{
			name: "TestFetchLogsFromServerWithSecondaryServer",
			args: args{
				base_path:        ".",
				file_name:        "testlog-not-from-primary",
				secondary_server: "true",
			},
			want:    []string{"this is from secondary server"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.secondary_server != "" {
				s := httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						type response struct {
							Logs     []string
							ErrorMsg string
						}

						resp := &response{
							Logs: []string{"this is from secondary server"},
						}
						json.NewEncoder(w).Encode(&resp)
					}),
				)
				defer s.Close()

				fmt.Println(s.URL)
				tt.args.secondary_server = s.URL
			}

			got, err := FetchLogsFromServer(tt.args.base_path, tt.args.file_name, tt.args.num_lines, tt.args.filter, tt.args.secondary_server)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchLogsFromServer() error = %v, wantErr %v", err, tt.wantErr)

				if tt.errMsg != err.Error() {
					t.Errorf("FetchLogsFromServer() error = %v, wantErr %v", err.Error(), tt.errMsg)
				}
				return
			}
			if !testEq(got, tt.want) {
				t.Errorf("FetchLogsFromServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
