package services

import "github.com/hugaojanuario/sentinel/internal/docker"

func ListContainers() (interface{}, error) {
	containers, err := docker.ListContainers()

	if err != nil {
		return nil, err
	}

	return containers, nil
}

func RestartContainer(id string) error {
	err := docker.RestartContainer(id)
	if err != nil {
		return err
	}

	return nil
}

func GetContainerLogs(id string) (string, error) {
	logs, err := docker.GetContainerLogs(id)
	if err != nil {
		return "", err
	}

	return logs, nil
}

func GetContainerStats(id string) (interface{}, error) {

	stats, err := docker.GetContainerStats(id)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
