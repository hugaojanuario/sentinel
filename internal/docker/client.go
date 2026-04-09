package docker

import (
	"sync"

	"github.com/docker/docker/client"
)

var (
	instance *client.Client
	once     sync.Once
)

func GetClient() (*client.Client, error) {
	var err error
	once.Do(func() {
		instance, err = client.NewClientWithOpts(client.FromEnv)
	})
	return instance, err
}
