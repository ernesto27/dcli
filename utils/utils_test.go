package utils

import "testing"

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
			want: "# title\n\n| column1 | column2 | column3 |\n| ------------- | ------------- | -------------|\n| row1 | row2 | row3 |\n| row4 | row5 | row6 |\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTable(tt.args.title, tt.args.columns, tt.args.rows); got != tt.want {
				t.Errorf("CreateTable() = %v, want %v", got, tt.want)
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
