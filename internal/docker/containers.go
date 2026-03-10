package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

type ContainerInfo struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Status string `json:"status"`
}

func ListContainers() ([]ContainerInfo, error){

	client, err := NewCLient()
	if err != nil{
		return nil, err
	}
	containers, err := client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil{
		return nil, err
	}
	
	var result []ContainerInfo
	for _, c := range containers{

		result = append(result, ContainerInfo{
			ID: c.ID,
			Name: c.Names[0],
			Image: c.Image,
			Status: c.Status,
		})
	}

	return result, nil

}

func RestartContainer(id string) error{

	client, err := NewCLient()
	if err != nil{
		return err
	}

	err = client.ContainerRestart(context.Background(), id, container.StopOptions{})

	if err != nil{
		return err
	}

	return nil

	
}