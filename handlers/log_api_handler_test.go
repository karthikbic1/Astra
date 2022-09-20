package handler

import (
	"net/url"
	"testing"
)

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
