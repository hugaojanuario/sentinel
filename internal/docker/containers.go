package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types/container"
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

	return client.ContainerRestart(context.Background(), id, container.StopOptions{})
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
