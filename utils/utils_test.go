package utils

import (
	"testing"

	"github.com/docker/docker/api/types/image"
)

func TestReverseSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    []image.HistoryResponseItem
		expected []image.HistoryResponseItem
	}{
		{
			name:     "empty slice",
			input:    []image.HistoryResponseItem{},
			expected: []image.HistoryResponseItem{},
		},
		{
			name: "slice with one element",
			input: []image.HistoryResponseItem{
				{
					ID: "123",
				},
			},
			expected: []image.HistoryResponseItem{
				{
					ID: "123",
				},
			},
		},
		{
			name: "slice with even number of elements",
			input: []image.HistoryResponseItem{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
				{
					ID: "789",
				},
				{
					ID: "012",
				},
			},
			expected: []image.HistoryResponseItem{
				{
					ID: "012",
				},
				{
					ID: "789",
				},
				{
					ID: "456",
				},
				{
					ID: "123",
				},
			},
		},
		{
			name: "slice with odd number of elements",
			input: []image.HistoryResponseItem{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
				{
					ID: "789",
				},
			},
			expected: []image.HistoryResponseItem{
				{
					ID: "789",
				},
				{
					ID: "456",
				},
				{
					ID: "123",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ReverseSlice(tc.input)

			if len(tc.input) != len(tc.expected) {
				t.Fatalf("expected length %d, but got %d", len(tc.expected), len(tc.input))
			}

			for i := range tc.input {
				if tc.input[i].ID != tc.expected[i].ID {
					t.Errorf("expected %v, but got %v", tc.expected[i], tc.input[i])
				}
			}
		})
	}
}

func TestCreateTable(t *testing.T) {
	type args struct {
		title   string
		columns []string
		rows    [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should create a table",
			args: args{
				title: "title",
				columns: []string{
					"column1",
					"column2",
					"column3",
				},
				rows: [][]string{
					{
						"row1",
						"row2",
						"row3",
					},
					{
						"row4",
						"row5",
						"row6",
					},
				},
			},
			want: "# title\n\n| column1 | column2 | column3 |\n| ------------- | ------------- | ------------- |\n| row1 | row2 | row3 |\n| row4 | row5 | row6 |\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTable(tt.args.title, tt.args.columns, tt.args.rows); got != tt.want {
				t.Errorf("CreateTable() = got %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetDockerHubURL(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should get official image format",
			args: args{
				image: "nginx",
			},
			want: "https://hub.docker.com/_/nginx",
		},
		{
			name: "should get other image format",
			args: args{
				image: "someuser/nameimage",
			},
			want: "https://hub.docker.com/r/someuser/nameimage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDockerHubURL(tt.args.image); got != tt.want {
				t.Errorf("GetDockerHubURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should format size bytes",
			args: args{
				size: 1,
			},
			want: "1 bytes",
		},
		{
			name: "should format size KB",
			args: args{
				size: 1024,
			},
			want: "1.00 KB",
		},
		{
			name: "should format size MB",
			args: args{
				size: 1024 * 1024,
			},
			want: "1.00 MB",
		},
		{
			name: "should format size GB",
			args: args{
				size: 1024 * 1024 * 1024,
			},
			want: "1.00 GB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSize(tt.args.size); got != tt.want {
				t.Errorf("FormatSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
