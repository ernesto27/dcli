package docker

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Docker struct {
	cli *client.Client
	ctx context.Context

	Containers []MyContainer
}

type MyNetwork struct {
	Name      string
	IPAddress string
	Gateway   string
}

type MyContainer struct {
	ID         string
	IDShort    string
	Name       string
	NameShort  string
	Image      string
	ImageShort string
	State      string
	Status     string
	Ports      []types.Port
	Size       string
	Command    string
	Env        []string
	Network    MyNetwork
}

func New(ctx context.Context) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Docker{
		cli: cli,
		ctx: ctx,
	}, nil
}

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

func (d *Docker) ContainerList() ([]MyContainer, error) {
	containers, err := d.cli.ContainerList(d.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	mc := []MyContainer{}
	for _, c := range containers {
		// print struct with nice format
		// fmt.Printf("%+v\n", c)

		cJSON, _, err := d.cli.ContainerInspectWithRaw(d.ctx, c.ID, true)

		if err != nil {
			continue
		}

		var name string
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}

		networkSettings := cJSON.NetworkSettings
		networkMode := string(cJSON.HostConfig.NetworkMode)

		mc = append(mc, MyContainer{
			ID:         c.ID,
			IDShort:    trimValue(c.ID, 10),
			Name:       name,
			NameShort:  trimValue(name, 20),
			Image:      c.Image,
			ImageShort: trimValue(c.Image, 20),
			State:      c.State,
			Status:     c.Status,
			Ports:      c.Ports,
			Size:       formatSize(*cJSON.SizeRootFs),
			Env:        cJSON.Config.Env,
			Network: MyNetwork{
				Name:      networkMode,
				IPAddress: networkSettings.IPAddress,
				Gateway:   networkSettings.Gateway,
			},
		})

		d.Containers = mc

	}
	return mc, nil
}

func (d *Docker) GetContainerByName(name string) (MyContainer, error) {
	for _, c := range d.Containers {
		if c.Name == name {
			return c, nil
		}
	}
	return MyContainer{}, fmt.Errorf("container %s not found", name)
}

func (d *Docker) ContainerRemove(containerID string) error {
	err := d.cli.ContainerRemove(d.ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	return err
}

func (d *Docker) ContainerStop(containerID string) error {
	timeout := 10
	err := d.cli.ContainerStop(d.ctx, containerID, container.StopOptions{
		Timeout: &timeout,
	})
	return err
}

func (d *Docker) ContainerStart(containerID string) error {
	err := d.cli.ContainerStart(d.ctx, containerID, types.ContainerStartOptions{})
	return err
}

func (d *Docker) ImageList() ([]types.ImageSummary, error) {
	images, err := d.cli.ImageList(d.ctx, types.ImageListOptions{})
	for index, image := range images {
		image.ID = strings.Replace(image.ID, "sha256:", "", -1)

		image.ID = trimValue(image.ID, 10)
		if len(image.RepoTags) > 0 {
			image.RepoTags[0] = trimValue(image.RepoTags[0], 40)
		}
		images[index] = image
	}

	return images, err
}

func (d *Docker) ImageRemove(imageID string) error {
	_, err := d.cli.ImageRemove(d.ctx, imageID, types.ImageRemoveOptions{
		PruneChildren: true,
		Force:         true,
	})
	return err
}

func (d *Docker) ServerVersion() (string, error) {
	typesVersion, err := d.cli.ServerVersion(d.ctx)
	if err != nil {
		return "", err
	}

	return typesVersion.Version, nil
}

func (d *Docker) ContainerRestart(containerID string) error {
	return d.cli.ContainerRestart(d.ctx, containerID, container.StopOptions{})
}

func (d *Docker) ContainerLogs(containerId string) (string, error) {
	out, err := d.cli.ContainerLogs(d.ctx, containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: true,
		Tail:       "800",
	})

	if err != nil {
		return "", err
	}

	defer out.Close()

	var response []string
	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		line := scanner.Text()
		response = append(response, line)
	}

	reverseLines(response)

	logs := ""
	for _, l := range response {
		logs += l[10:] + "\n"
	}

	return logs, nil
}

func reverseLines(lines []string) {
	i := 0
	j := len(lines) - 1

	for i < j {
		lines[i], lines[j] = lines[j], lines[i]
		i++
		j--
	}
}

func trimValue(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
