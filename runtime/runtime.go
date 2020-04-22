// Runtimes for different programming languages
// For now I'm going to spawn one container per runtime
package runtime

import(
  "context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Runtime struct{
  Name string
  Exec string
  Tag string
  Image string
  ExistsLocally *bool
}

func(r Runtime) CheckIfExistsLocally(cli *client.Client, ctx context.Context) bool{

  if r.ExistsLocally != nil{
    return *r.ExistsLocally
  }else{
    exists := true
    AllImagesInHost, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
    if err != nil {
      panic(err)
    }

    for _, image := range AllImagesInHost{
        for _, tag := range image.RepoTags{
          if r.Tag == tag{
            r.ExistsLocally = &exists
            break
          }
        }

        if r.ExistsLocally != nil{
          if *r.ExistsLocally == true{
            break
          }
        }
    }

    if r.ExistsLocally == nil{
      exists := false
      r.ExistsLocally = &exists
    }
  }

  return *r.ExistsLocally

}

var Runtimes = []Runtime {
  Runtime{Name: "theangrydev_python", Tag: "python:latest", Exec: "python", Image: "docker.io/python"},
  Runtime{Name: "theangrydev_ruby", Tag: "ruby:latest", Exec: "ruby", Image: "docker.io/ruby"},
  Runtime{Name: "theangrydev_node", Tag: "node:latest", Exec: "node", Image: "docker.io/node"}}
