package docker

import (
	"github.com/docker/docker/client"
)

func NewCLient() (*client.Client, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return cli, nil
}
