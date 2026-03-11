package docker

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types/container"
)

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

func ListContainers() ([]ContainerInfo, error) {

	client, err := NewCLient()
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

	client, err := NewCLient()
	if err != nil {
		return err
	}

	return client.ContainerRestart(context.Background(), id, container.StopOptions{})
}

func GetContainerLogs(id string) (string, error) {
	client, err := NewCLient()
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

func GetContainerStats(id string) (interface{}, error) {

	client, err := NewCLient()
	if err != nil {
		return nil, err
	}

	stats, err := client.ContainerStats(context.Background(), id, false)
	if err != nil {
		return nil, err
	}

	defer stats.Body.Close()

	var data map[string]interface{}

	err = json.NewDecoder(stats.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
