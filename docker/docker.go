package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ernesto27/dcli/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Docker struct {
	cli *client.Client
	ctx context.Context

	Containers []MyContainer
	Images     []MyImage
	Networks   []MyNetwork
	Volumes    []MyVolume
	Stacks     []MyStack
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
	ID           string
	IDShort      string
	Name         string
	NameShort    string
	Image        string
	ImageShort   string
	State        string
	Status       string
	Ports        []types.Port
	Size         string
	SizeOriginal int64
	Command      string
	Env          []string
	ReadOnly     bool
	MountedAt    string
	Network      MyNetwork
	Mounts       []types.MountPoint
}

type MyContainerStats struct {
	ID       string
	CPUPer   float64
	MemUsage string
	MemLimit string
	MemPer   float64
	PID      uint64
}

type MyImage struct {
	Summary types.ImageSummary
	Inspect types.ImageInspect
	History []image.HistoryResponseItem
}

type MyVolume struct {
	Volume     *volume.Volume
	Containers []MyContainer
}

type MyStack struct {
	MyNetwork
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
	return utils.FormatSize(i.Summary.Size)
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

		readOnly := false
		mountedAt := ""

		if len(c.Mounts) > 0 {
			readOnly = c.Mounts[0].RW == false
			mountedAt = c.Mounts[0].Destination
		}

		mc = append(mc, MyContainer{
			ID:           c.ID,
			IDShort:      utils.TrimValue(c.ID, 10),
			Name:         name,
			NameShort:    utils.TrimValue(name, 20),
			Image:        c.Image,
			ImageShort:   utils.TrimValue(c.Image, 20),
			State:        c.State,
			Status:       c.Status,
			Ports:        c.Ports,
			Size:         utils.FormatSize(*cJSON.SizeRootFs),
			SizeOriginal: *cJSON.SizeRootFs,
			Env:          cJSON.Config.Env,
			Command:      strings.Join(cJSON.Config.Entrypoint, " ") + " " + strings.Join(cJSON.Config.Cmd, " "),
			ReadOnly:     readOnly,
			MountedAt:    mountedAt,
			Network: MyNetwork{
				Name:      networkMode,
				IPAddress: ipAddress,
				Gateway:   gateway,
			},
			Mounts: c.Mounts,
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
	return d.cli.ContainerStart(d.ctx, containerID, types.ContainerStartOptions{})
}

func (d *Docker) ContainerPause(containerID string) error {
	return d.cli.ContainerPause(d.ctx, containerID)
}

func (d *Docker) ContainerUnpause(containerID string) error {
	return d.cli.ContainerUnpause(d.ctx, containerID)
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

func (d *Docker) ImageRemove(imageID string, force bool) error {
	_, err := d.cli.ImageRemove(d.ctx, imageID, types.ImageRemoveOptions{
		PruneChildren: true,
		Force:         force,
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
		if len(l) > 10 {
			logs += l[10:] + "\n"
		}
	}

	return logs, nil
}

func (d *Docker) ContainerStats(containerID string) (MyContainerStats, error) {
	s, err := d.cli.ContainerStats(d.ctx, containerID, false)
	if err != nil {
		panic(err)
	}
	defer s.Body.Close()

	var containerStats types.Stats
	dec := json.NewDecoder(s.Body)
	if err := dec.Decode(&containerStats); err != nil {
		panic(err)
	}

	cpuPercentage := calculateCPUPercentage(&containerStats)
	memUsage, memLimit := calculateMemoryUsage(&containerStats)
	memPercentage := calculateMemoryPercentage(memUsage, memLimit)

	cs := MyContainerStats{
		ID:       containerID,
		MemUsage: formatSizeStats(memUsage),
		MemLimit: formatSizeStats(memLimit),
		MemPer:   memPercentage,
		CPUPer:   cpuPercentage,
		PID:      containerStats.PidsStats.Current,
	}

	return cs, err
}

func formatSizeStats(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}

	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f%s", size, units[unitIndex])
}

func calculateCPUPercentage(stats *types.Stats) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

	cpuPercentage := 0.0
	if systemDelta > 0.0 {
		cpuPercentage = (cpuDelta / systemDelta) * 100.0
	}

	return cpuPercentage
}

func calculateMemoryUsage(stats *types.Stats) (float64, float64) {
	memUsage := float64(stats.MemoryStats.Usage)
	memLimit := float64(stats.MemoryStats.Limit)
	return memUsage, memLimit
}

func calculateMemoryPercentage(memUsage, memLimit float64) float64 {
	if memLimit <= 0.0 {
		return 0.0
	}

	return (memUsage / memLimit) * 100.0
}

type NetworkIO struct {
	Input  float64
	Output float64
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
		for _, c := range d.Containers {
			if c.Network.Name == network.Name {
				containers = append(containers, c)
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

func (d *Docker) NetworkRemove(networkID string) error {
	return d.cli.NetworkRemove(d.ctx, networkID)
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

func (d *Docker) VolumeList() ([]MyVolume, error) {
	options := volume.ListOptions{}
	vl, err := d.cli.VolumeList(d.ctx, options)
	if err != nil {
		return []MyVolume{}, err
	}

	mvol := []MyVolume{}
	for _, v := range vl.Volumes {
		name := v.Name

		containers := []MyContainer{}
		for _, c := range d.Containers {
			for _, mount := range c.Mounts {
				if mount.Type == "volume" && mount.Name == name {
					containers = append(containers, c)
				}
			}
		}

		mvol = append(mvol, MyVolume{
			Volume:     v,
			Containers: containers,
		})
	}

	d.Volumes = mvol
	return mvol, nil
}

func (d *Docker) GetVolumeByName(name string) (MyVolume, error) {
	for _, v := range d.Volumes {
		if v.Volume.Name == name {
			return v, nil
		}
	}

	return MyVolume{}, fmt.Errorf("volume %s not found", name)
}

func (d *Docker) VolumeRemove(volumeID string) error {
	return d.cli.VolumeRemove(d.ctx, volumeID, false)
}

func (d *Docker) StackList() ([]MyStack, error) {
	stacks := []MyStack{}
	networks, err := d.NetworkList()
	if err != nil {
		return stacks, err
	}

	for _, n := range networks {
		if n.Resource.Labels["com.docker.compose.project"] != "" {
			stacks = append(stacks, MyStack{n})
		}
	}

	d.Stacks = stacks
	return stacks, nil
}

func (d *Docker) GetStackByName(name string) (MyStack, error) {

	for _, s := range d.Stacks {
		if s.Resource.Name == name {
			return s, nil
		}
	}

	return MyStack{}, fmt.Errorf("stack %s not found", name)
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

func (d *Docker) GetAllImagesSize() string {
	var size int64
	for _, image := range d.Images {
		size += image.Summary.Size
	}

	return utils.FormatSize(size)
}

func (d *Docker) GetAllContainersSize() string {
	var size int64

	for _, container := range d.Containers {
		size += container.SizeOriginal
	}

	return utils.FormatSize(size)
}

func (d *Docker) GetImagesOrderBySize(desc bool) []MyImage {
	sortedImages := make([]MyImage, len(d.Images))
	copy(sortedImages, d.Images)
	sort.Slice(sortedImages, func(i, j int) bool {
		if desc {
			return sortedImages[i].Summary.Size > sortedImages[j].Summary.Size
		} else {
			return sortedImages[i].Summary.Size < sortedImages[j].Summary.Size
		}
	})
	return sortedImages
}

func (d *Docker) GetContainersOrderBySize(desc bool) []MyContainer {
	sortedContainers := make([]MyContainer, len(d.Containers))
	copy(sortedContainers, d.Containers)
	sort.Slice(sortedContainers, func(i, j int) bool {
		if desc {
			return sortedContainers[i].SizeOriginal > sortedContainers[j].SizeOriginal
		} else {
			return sortedContainers[i].SizeOriginal < sortedContainers[j].SizeOriginal
		}
	})
	return sortedContainers
}
