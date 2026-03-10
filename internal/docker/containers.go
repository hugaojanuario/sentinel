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

	cli, err := NewCLient()
	if err != nil{
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
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