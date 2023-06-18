package utils

import "testing"

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
			name: "official image format",
			args: args{
				image: "nginx",
			},
			want: "https://hub.docker.com/_/nginx",
		},
		{
			name: "Other image format",
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
