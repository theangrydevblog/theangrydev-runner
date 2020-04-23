package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"./runtime"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var ctx context.Context = context.Background()
var wg sync.WaitGroup

func createContainer(cli *client.Client, r runtime.Runtime) {

	if r.CheckIfExistsLocally(ctx, cli) != true {
		fmt.Printf("%s not found in host. Pulling from container registry...\n", r.Tag)
		out, pullError := cli.ImagePull(ctx, r.Tag, types.ImagePullOptions{All: false})
		if pullError != nil {
			panic(pullError)
		}

		defer out.Close()
		io.Copy(os.Stdout, out)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: r.Image,
		Cmd:   []string{"/bin/bash"},
		Tty:   true,
	}, nil, nil, r.Name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Created container %s : %s\n", r.Name, resp.ID)
	wg.Done()
}

func runCode(r runtime.Runtime, code string) {

}

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	wg.Add(len(runtime.Runtimes))

	for _, r := range runtime.Runtimes {
		go createContainer(cli, r)
	}

	wg.Wait()

}
