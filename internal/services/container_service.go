package services

import "github.com/hugaojanuario/sentinel/internal/docker"

func ListContainers() (interface{}, error) {
	containers, err := docker.ListContainers()

	if err != nil{
		return nil, err
	}

	return containers, nil
}