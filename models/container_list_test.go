package models

import (
	"reflect"
	"testing"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/table"
)

func TestGetContainerRows(t *testing.T) {
	type args struct {
		containers []docker.MyContainer
		query      string
	}
	tests := []struct {
		name string
		args args
		want []table.Row
	}{
		{
			name: "should get all results if query is empty",
			args: args{
				containers: []docker.MyContainer{
					{
						ID:    "1234567890",
						Name:  "test",
						Image: "test",
						State: running,
					},
					{
						ID:    "12345678902",
						Name:  "test2",
						Image: "test2",
						State: exited,
					},
				},
				query: "",
			},
			want: []table.Row{{"1234567890", "test", "test", "", "", "\033[32m\u2191\033[0m " + running}, {"12345678902", "test2", "test2", "", "", "\033[31m\u2193\033[0m " + exited}},
		},
		{
			name: "should get filtered results when search by name",
			args: args{
				containers: []docker.MyContainer{
					{
						ID:    "1234567890",
						Name:  "nginx",
						Image: "test",
						State: running,
					},
					{
						ID:    "1234567890",
						Name:  "test",
						Image: "test",
						State: running,
					},
				},
				query: "nginx",
			},
			want: []table.Row{{"1234567890", "nginx", "test", "", "", "\033[32m\u2191\033[0m " + running}},
		},
		{
			name: "should get filtered results when search by image",
			args: args{
				containers: []docker.MyContainer{
					{
						ID:    "1234567890",
						Name:  "test",
						Image: "mysql",
						State: running,
					},
					{
						ID:    "1234567890",
						Name:  "test",
						Image: "test",
						State: running,
					},
				},
				query: "mysq",
			},
			want: []table.Row{{"1234567890", "test", "mysql", "", "", "\033[32m\u2191\033[0m " + running}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContainerRows(tt.args.containers, tt.args.query)

			if len(got) != len(tt.want) {
				t.Errorf("GetContainerRows() got = %v, want %v", len(got), len(tt.want))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContainerRows() got = %v, want %v", got, tt.want)
			}

		})
	}

}
