package docker

import (
	"bufio"
	"context"
	"dockerniceui/utils"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Docker struct {
	cli *client.Client
	ctx context.Context

	Containers []MyContainer
	Images     []MyImage
	Networks   []MyNetwork
}

type MyNetwork struct {
	Name       string
	IPAddress  string
	Gateway    string
	Subnet     string
	Resource   types.NetworkResource
	Containers []MyContainer
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

type MyImage struct {
	Summary types.ImageSummary
	Inspect types.ImageInspect
	History []image.HistoryResponseItem
}

func (i *MyImage) GetFormatTimestamp() string {
	if i.Summary.Created == 0 {
		return ""
	}
	return FormatTimestamp(i.Summary.Created)
}

func (i *MyImage) GetFormatSize() string {
	if i.Summary.Size == 0 {
		return ""
	}
	return formatSize(i.Summary.Size)
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

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	duration := time.Since(t)

	months := int(duration.Hours() / (24 * 30))
	weeks := int(duration.Hours() / (24 * 7))
	days := int(duration.Hours() / 24)

	plural := "s"
	if months > 0 {
		if months == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d month%s ago", months, plural)
	} else if weeks > 0 {
		if weeks == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d week%s ago", weeks, plural)
	} else if days > 0 {
		if days == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d day%s ago", days, plural)
	} else {
		return "today"
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
		cJSON, _, err := d.cli.ContainerInspectWithRaw(d.ctx, c.ID, true)

		if err != nil {
			continue
		}

		name := d.getContainerName(c.Names)
		networkSettings := cJSON.NetworkSettings
		networkMode := string(cJSON.HostConfig.NetworkMode)

		defaultNetwork := "default"
		ipAddress := d.getContainerIP(cJSON)

		gateway := networkSettings.Gateway
		if networkMode != defaultNetwork {
			gateway = networkSettings.Networks[networkMode].Gateway
		}

		mc = append(mc, MyContainer{
			ID:         c.ID,
			IDShort:    utils.TrimValue(c.ID, 10),
			Name:       name,
			NameShort:  utils.TrimValue(name, 20),
			Image:      c.Image,
			ImageShort: utils.TrimValue(c.Image, 20),
			State:      c.State,
			Status:     c.Status,
			Ports:      c.Ports,
			Size:       formatSize(*cJSON.SizeRootFs),
			Env:        cJSON.Config.Env,
			Command:    strings.Join(cJSON.Config.Entrypoint, " ") + " " + strings.Join(cJSON.Config.Cmd, " "),
			Network: MyNetwork{
				Name:      networkMode,
				IPAddress: ipAddress,
				Gateway:   gateway,
			},
		})

		d.Containers = mc

	}
	return mc, nil
}

func (d *Docker) getContainerName(names []string) string {
	if len(names) > 0 {
		return names[0][1:]
	}
	return ""
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

func (d *Docker) ImageList() ([]MyImage, error) {
	images, err := d.cli.ImageList(d.ctx, types.ImageListOptions{})
	myImages := []MyImage{}

	for index, image := range images {
		imageInspect, _, err := d.cli.ImageInspectWithRaw(d.ctx, image.ID)
		if err != nil {
			fmt.Println(err)
		}

		fullID := strings.Replace(image.ID, "sha256:", "", -1)

		image.ID = utils.TrimValue(fullID, 10)
		if len(image.RepoTags) > 0 {
			image.RepoTags[0] = utils.TrimValue(image.RepoTags[0], 40)
		} else {
			image.RepoTags = []string{"<none>"}
		}
		images[index] = image

		h, err := d.cli.ImageHistory(d.ctx, fullID)
		if err != nil {
			fmt.Println(err)
		}

		utils.ReverseSlice(h)

		myImages = append(myImages, MyImage{
			Summary: image,
			Inspect: imageInspect,
			History: h,
		})
	}

	d.Images = myImages
	return myImages, err
}

func (d *Docker) GetImageByID(ID string) (MyImage, error) {
	for _, i := range d.Images {
		if i.Summary.ID == ID {
			return i, nil
		}
	}
	return MyImage{}, fmt.Errorf("image %s not found", ID)
}

func (d *Docker) ImageRemove(imageID string) error {
	_, err := d.cli.ImageRemove(d.ctx, imageID, types.ImageRemoveOptions{
		PruneChildren: true,
		Force:         false,
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

	utils.ReverseLines(response)

	logs := ""
	for _, l := range response {
		logs += l[10:] + "\n"
	}

	return logs, nil
}

func (d *Docker) NetworkList() ([]MyNetwork, error) {
	myNetwork := []MyNetwork{}
	networks, err := d.cli.NetworkList(d.ctx, types.NetworkListOptions{})
	if err != nil {
		return myNetwork, err
	}

	for _, n := range networks {
		subnet := ""
		gateway := ""
		if len(n.IPAM.Config) > 0 {
			subnet = n.IPAM.Config[0].Subnet
			gateway = n.IPAM.Config[0].Gateway
		}

		network, err := d.cli.NetworkInspect(context.Background(), n.ID, types.NetworkInspectOptions{})
		if err != nil {
			fmt.Println(err)
		}

		containers := []MyContainer{}
		for containerID := range network.Containers {
			container, _, err := d.cli.ContainerInspectWithRaw(context.Background(), containerID, true)
			if err == nil {
				n := MyNetwork{
					IPAddress: d.getContainerIP(container),
				}
				containers = append(containers, MyContainer{
					Name:    d.getContainerName([]string{container.Name}),
					Network: n,
				})
			}
		}

		myNetwork = append(myNetwork, MyNetwork{
			Resource:   network,
			Gateway:    gateway,
			Subnet:     subnet,
			Containers: containers,
		})
	}

	d.Networks = myNetwork

	return myNetwork, nil
}

func (d *Docker) getContainerIP(c types.ContainerJSON) string {
	networkSettings := c.NetworkSettings
	networkMode := string(c.HostConfig.NetworkMode)

	defaultNetwork := "default"
	ipAddress := networkSettings.IPAddress
	if networkMode != defaultNetwork {
		ipAddress = networkSettings.Networks[networkMode].IPAddress
	}
	return ipAddress
}

func (d *Docker) GetNetworkByName(name string) (MyNetwork, error) {
	for _, n := range d.Networks {
		if n.Resource.Name == name {
			return n, nil
		}
	}
	return MyNetwork{}, fmt.Errorf("network %s not found", name)
}

func (d *Docker) Events() {
	go func() {
		eventStream, err := d.cli.Events(d.ctx, types.EventsOptions{})
		for {
			select {
			case err := <-err:
				fmt.Println(err)

			case <-eventStream:
				d.ContainerList()
			}
		}
	}()
}
