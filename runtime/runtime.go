// Package runtime provides structures to deal with different runtime environemnts
package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Runtime represents a runtime environment/Docker container(s). Eg: Node.js, CPython, CRuby etc
type Runtime struct {
	Name          string
	Exec          []string
	Tag           string
	Image         string
	ID            *string
	ExistsLocally *bool
}

// ExecuteCode runs the source code in the relevant runtime environment
func (r *Runtime) ExecuteCode(ctx context.Context, cli *client.Client, code string) string {

	execConfig := types.ExecConfig{
		Cmd:          append(r.Exec, code),
		AttachStdout: true,
		AttachStderr: true,
	}
	start, err := cli.ContainerExecCreate(ctx, *r.ID, execConfig)
	if err != nil {
		panic(err)
	}

	config := types.ExecStartCheck{}
	res, er := cli.ContainerExecAttach(ctx, start.ID, config)
	if er != nil {
		panic(er)
	}

	execErr := cli.ContainerExecStart(context.Background(), start.ID, types.ExecStartCheck{
		Tty: true,
	})

	if execErr != nil {
		panic(execErr)
	}

	// inspect, inspectErr := cli.ContainerExecInspect(context.Background(), start.ID)
	// if inspectErr != nil {
	// 	panic(inspectErr)
	// }

	content, _, _ := res.Reader.ReadLine()
	// io.Copy(os.Stdout, content)
	return string(content)

}

// CheckIfExistsLocally checks if the desired image is already present in the host machine
func (r *Runtime) CheckIfExistsLocally(ctx context.Context, cli *client.Client) bool {

	if r.ExistsLocally != nil {
		return *r.ExistsLocally
	}
	exists := true
	AllImagesInHost, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, image := range AllImagesInHost {
		for _, tag := range image.RepoTags {
			if r.Tag == tag {
				r.ExistsLocally = &exists
				break
			}
		}

		if r.ExistsLocally != nil {
			if *r.ExistsLocally == true {
				break
			}
		}
	}

	if r.ExistsLocally == nil {
		exists := false
		r.ExistsLocally = &exists
	}

	return *r.ExistsLocally

}

// UpdateID checks if there is a container ID attached to this runtime
// If not, then either 1) There are no containers running 2) Containers are running but weren't spawned by the current main process
// If 2) is the case, grab the ID of the currently active container and store it in the struct
func (r *Runtime) UpdateID(ctx context.Context, cli *client.Client) *string {
	if r.ID == nil {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			for _, name := range container.Names {
				if "/"+r.Name == name {
					r.ID = &container.ID
					break
				}
			}

			if r.ID != nil {
				fmt.Printf("Container %s found for runtime %s\n", *r.ID, r.Name)
				break
			}
		}

	}

	return r.ID

}

// Runtimes is a list of runtime environments we want to spawn to service requests
// Currently only 1 runtime per language should suffice
var Runtimes = []Runtime{
	Runtime{Name: "theangrydev_python", Tag: "python:latest", Exec: []string{"python", "-c"}, Image: "docker.io/python"},
	Runtime{Name: "theangrydev_ruby", Tag: "ruby:latest", Exec: []string{"ruby", "-e"}, Image: "docker.io/ruby"},
	Runtime{Name: "theangrydev_rust", Tag: "rust:latest", Exec: []string{}, Image: "docker.io/rust"},
	Runtime{Name: "theangrydev_node", Tag: "node:latest", Exec: []string{}, Image: "docker.io/node"}}
