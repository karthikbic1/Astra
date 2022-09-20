package core

import "testing"

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
		want    string
		wantErr bool
	}{
		{
			name: "TestFetchLogsFromServer",
			args: args{
				base_path: ".",
				file_name: "testlog",
			},
			want: `
this is 9 line.
this is 8 line 
this is 7 line
this is 6 line.
this is 5 line 今日は
this is 4 line
this is 3 line.
this is 2 line
this is 1 line
`,
			wantErr: false,
		},
		{
			name: "TestFetchLogsFromServerWithNumLines",
			args: args{
				base_path: ".",
				file_name: "testlog",
				num_lines: 2,
			},
			want: `
this is 9 line.
this is 8 line `,
			wantErr: false,
		},

		{
			name: "TestFetchLogsFromServerWithFilter",
			args: args{
				base_path: ".",
				file_name: "testlog",
				filter:    "5 line",
			},
			want: `
this is 5 line 今日は`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchLogsFromServer(tt.args.base_path, tt.args.file_name, tt.args.num_lines, tt.args.filter, tt.args.secondary_server)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchLogsFromServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchLogsFromServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
