package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrContainersNotRunning = errors.New("containers not running")
)

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type ContainerStats struct {
	CpuStats    CpuStats    `json:"cpu_stats"`
	PreCpuStats CpuStats    `json:"precpu_stats"`
	MemoryStats MemoryStats `json:"memory_stats"`
}

type CpuStats struct {
	CpuUsage       CpuUsage `json:"cpu_usage"`
	SystemCpuUsage uint64   `json:"system_cpu_usage"`
	OnlineCpus     uint64   `json:"online_cpus"`
}

type CpuUsage struct {
	TotalUsage uint64 `json:"total_usage"`
}

type MemoryStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

func ListContainers() ([]ContainerInfo, error) {

	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	containers, err := client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo
	for _, c := range containers {

		result = append(result, ContainerInfo{
			ID:     c.ID,
			Name:   c.Names[0],
			Image:  c.Image,
			Status: c.Status,
		})
	}

	return result, nil

}

func RestartContainer(id string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	err = client.ContainerRestart(context.Background(), id, container.StopOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "No such container") {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func GetContainerLogs(id string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	menu := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "50",
	}

	reader, err := client.ContainerLogs(context.Background(), id, menu)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetContainerStats(id string) (ContainerStats, error) {

	client, err := GetClient()
	if err != nil {
		return ContainerStats{}, err
	}

	stats, err := client.ContainerStats(context.Background(), id, false)
	if err != nil {
		return ContainerStats{}, err
	}

	defer stats.Body.Close()

	var containerStats ContainerStats

	err = json.NewDecoder(stats.Body).Decode(&containerStats)
	if err != nil {
		return ContainerStats{}, err
	}

	return containerStats, nil
}

func StreamContainerLogs(id string) (io.ReadCloser, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "50",
	}

	reader, err := client.ContainerLogs(context.Background(), id, options)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

type ContainerState struct {
	ID         string
	Name       string
	Status     string
	OOMKilled  bool
	ExitCode   int
	ExitError  string
	FinishedAt time.Time
}

func ListAllContainers() ([]ContainerInfo, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		result = append(result, ContainerInfo{
			ID:     c.ID,
			Name:   name,
			Image:  c.Image,
			Status: c.Status,
		})
	}

	return result, nil
}

func InspectContainer(id string) (ContainerState, error) {
	cl, err := GetClient()
	if err != nil {
		return ContainerState{}, err
	}

	info, err := cl.ContainerInspect(context.Background(), id)
	if err != nil {
		return ContainerState{}, err
	}

	name := strings.TrimPrefix(info.Name, "/")

	var finishedAt time.Time
	if t, err := time.Parse(time.RFC3339Nano, info.State.FinishedAt); err == nil {
		finishedAt = t
	}

	return ContainerState{
		ID:         info.ID,
		Name:       name,
		Status:     info.State.Status,
		OOMKilled:  info.State.OOMKilled,
		ExitCode:   info.State.ExitCode,
		ExitError:  info.State.Error,
		FinishedAt: finishedAt,
	}, nil
}
