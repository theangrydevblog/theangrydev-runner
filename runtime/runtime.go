// Package runtime provides structures to deal with different runtime environemnts
package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Runtime represents a runtime environment. Eg: Node.js, CPython, CRuby etc
type Runtime struct {
	Name          string
	Exec          string
	Tag           string
	Image         string
	ExistsLocally *bool
}

// CheckIfExistsLocally checks if the desired image is already present in the host machine
func (r Runtime) CheckIfExistsLocally(ctx context.Context, cli *client.Client) bool {

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

// Runtimes is a list of runtime environments we want to spawn to service requests
// Currently only 1 runtime per language should suffice
var Runtimes = []Runtime{
	Runtime{Name: "theangrydev_python", Tag: "python:latest", Exec: "python", Image: "docker.io/python"},
	Runtime{Name: "theangrydev_ruby", Tag: "ruby:latest", Exec: "ruby", Image: "docker.io/ruby"},
	Runtime{Name: "theangrydev_rust", Tag: "rust:latest", Exec: "ruby", Image: "docker.io/rust"},
	Runtime{Name: "theangrydev_node", Tag: "node:latest", Exec: "node", Image: "docker.io/node"}}
