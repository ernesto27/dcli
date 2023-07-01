package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

func TestGetContainerIP(t *testing.T) {
	type args struct {
		container types.ContainerJSON
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should get container ip",
			args: args{
				container: types.ContainerJSON{
					ContainerJSONBase: &types.ContainerJSONBase{
						HostConfig: &container.HostConfig{
							NetworkMode: "bridge",
						},
					},
					NetworkSettings: &types.NetworkSettings{
						Networks: map[string]*network.EndpointSettings{
							"bridge": {
								IPAddress: "172.0.0.10",
							},
						},
					},
				},
			},
			want: "172.0.0.10",
		},
		{
			name: "should get empty string if network mode is default",
			args: args{
				container: types.ContainerJSON{
					ContainerJSONBase: &types.ContainerJSONBase{
						HostConfig: &container.HostConfig{
							NetworkMode: "default",
						},
					},
					NetworkSettings: &types.NetworkSettings{
						Networks: map[string]*network.EndpointSettings{
							"default": {
								IPAddress: "",
							},
						},
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Docker{}
			if got := d.getContainerIP(tt.args.container); got != tt.want {
				t.Errorf("getContainerIP() = got%v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContainerName(t *testing.T) {
	type args struct {
		names []string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should get container name without slash",
			args: args{
				names: []string{"/container_name"},
			},
			want: "container_name",
		},
		{
			name: "should get empty string args is empty",
			args: args{
				names: []string{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Docker{}
			if got := d.getContainerName(tt.args.names); got != tt.want {
				t.Errorf("getContainerName() = got %v, want %v", got, tt.want)
			}
		})
	}
}
