package main

import(
  "fmt"
  "context"
  "sync"
  "./runtime"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)


var ctx context.Context = context.Background()
var wg sync.WaitGroup

func createContainer(cli *client.Client, name string, image string){
  resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image: image,
    Cmd: []string{"/bin/bash"},
    Tty: true,
  }, nil, nil, name)
  if err != nil {
    // TODO: If image not found, Pull it from Docker hub
    panic(err)
  }

  if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    panic(err)
  }

  fmt.Printf("Created container %s : %s\n", name, resp.ID)
  wg.Done()
}


func destroyContainer(cli *client.Client, name string, image string){

}


func runCode(r runtime.Runtime, code string){

}

func main(){

  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    panic(err)
  }

  wg.Add(len(runtime.Runtimes))

  for _, r := range runtime.Runtimes{
    go createContainer(cli, r.Name, r.Image)
  }

  wg.Wait()

}
