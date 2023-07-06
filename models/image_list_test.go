package models

import (
	"reflect"
	"testing"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/table"
	"github.com/docker/docker/api/types"
)

func TestGetImagesRows(t *testing.T) {
	type args struct {
		images []docker.MyImage
		query  string
	}
	tests := []struct {
		name string
		args args
		want []table.Row
	}{
		{
			name: "should get all results if query is empty",
			args: args{
				images: []docker.MyImage{
					{
						Summary: types.ImageSummary{
							ID:       "1234567890",
							RepoTags: []string{"test:test"},
						},
					},
					{
						Summary: types.ImageSummary{
							ID:       "12345678902",
							RepoTags: []string{"test2:test2"},
						},
					},
				},
				query: "",
			},
			want: []table.Row{{"1234567890", "test:test", "", ""}, {"12345678902", "test2:test2", "", ""}},
		},
		{
			name: "should get filtered results when search by name",
			args: args{
				images: []docker.MyImage{
					{
						Summary: types.ImageSummary{
							ID:       "1234567890",
							RepoTags: []string{"nginx:latest"},
						},
					},
					{
						Summary: types.ImageSummary{
							ID:       "12345678902",
							RepoTags: []string{"test2:test2"},
						},
					},
				},
				query: "ngin",
			},
			want: []table.Row{{"1234567890", "nginx:latest", "", ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageRows(tt.args.images, tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetImagesRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
