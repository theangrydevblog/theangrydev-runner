// Package runtime provides structures to deal with different runtime environemnts
package runtime

import (
	"context"

	"../container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Runtime represents a runtime environment/Docker container(s). Eg: Node.js, CPython, CRuby etc
type Runtime struct {
	Name          string
	Exec          []string
	Tag           string
	Image         string
	ExistsLocally *bool
	Container     *container.Container
}

// ExecuteCode runs the source code in the relevant runtime environment
func (r *Runtime) ExecuteCode(ctx context.Context, cli *client.Client, code string) string {

	execConfig := types.ExecConfig{
		Cmd:          append(r.Exec, code),
		AttachStdout: true,
		AttachStderr: true,
	}
	start, err := cli.ContainerExecCreate(ctx, *r.Container.ID, execConfig)
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

// Runtimes is a map of runtime environments we want to spawn to service requests
var Runtimes = map[string]*Runtime{
	"Python": &Runtime{
		Name:      "python",
		Tag:       "python:latest",
		Exec:      []string{"python", "-c"},
		Container: &container.Container{Name: "theangrydev_python"},
		Image:     "docker.io/python"},
	"Ruby": &Runtime{
		Name:      "ruby",
		Tag:       "ruby:latest",
		Exec:      []string{"ruby", "-e"},
		Container: &container.Container{Name: "theangrydev_ruby"},
		Image:     "docker.io/ruby"},
	"Rust": &Runtime{
		Name:      "rust",
		Tag:       "rust:latest",
		Exec:      []string{},
		Container: &container.Container{Name: "theangrydev_rust"},
		Image:     "docker.io/rust"},
	"JavaScript": &Runtime{
		Name:      "node",
		Tag:       "node:latest",
		Exec:      []string{},
		Container: &container.Container{Name: "theangrydev_node"},
		Image:     "docker.io/node"}}
