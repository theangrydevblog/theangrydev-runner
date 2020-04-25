package container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Container represents a Docker container. It is responsible for making sure that a Runtime always has a Docker container attached to it
type Container struct {
	Name string
	ID   *string
}

// Health is a wrapper for health check messages. Should be sent to the alertChannel
type Health struct {
	MessageCode bool
	Container   *Container
}

// CheckIfRunning checks if container is running or stopped
func (c *Container) CheckIfRunning(ctx context.Context, cli *client.Client) bool {
	runningContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, runningContainer := range runningContainers {
		for _, containerName := range runningContainer.Names {
			if "/"+c.Name == containerName {
				return true
			}
		}
	}

	return false

}

// UpdateID grabs the container ID of the container and updates the struct. Only run this if a container is already running
func (c *Container) UpdateID(ctx context.Context, cli *client.Client) {
	runningContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, runningContainer := range runningContainers {
		for _, containerName := range runningContainer.Names {
			if "/"+c.Name == containerName {
				c.ID = &runningContainer.ID
				break
			}
		}

		if c.ID != nil {
			fmt.Printf("Updated container %s's ID to : %s\n", c.Name, *c.ID)
			break
		}
	}

}

// RestartContainer restarts a shut down container
func (c *Container) RestartContainer(ctx context.Context, cli *client.Client) {
	fmt.Printf("Restarting container %s (%s)\n", c.Name, *c.ID)
	if err := cli.ContainerStart(ctx, *c.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Restarted container %s\n", c.Name)
}

// StartContainer creates a container and starts it
func (c *Container) StartContainer(ctx context.Context, cli *client.Client, image string) {
	fmt.Printf("Creating container %s\n", c.Name)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash"},
		Tty:   true,
	}, nil, nil, c.Name)
	if err != nil {
		panic(err)
	}

	c.ID = &resp.ID

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Started container %s\n", c.Name)

}

// Monitor checks the health of a container every 2 seconds and sends a message in the alertChannel
func (c *Container) Monitor(ctx context.Context, cli *client.Client, alertChannel chan Health) {
	for {
		time.Sleep(2 * time.Second)
		if c.CheckIfRunning(ctx, cli) == true {
			alertChannel <- Health{MessageCode: true, Container: c}
		} else {
			alertChannel <- Health{MessageCode: false, Container: c}
		}
	}
}
