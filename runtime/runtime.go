// Runtimes for different programming languages
// For now I'm going to spawn one container per runtime
package runtime

type Runtime struct{
  Name string
  Exec string
  Image string
}

var Runtimes = []Runtime {
  Runtime{Name: "python", Exec: "python", Image: "docker.io/python"},
  Runtime{Name: "ruby", Exec: "ruby", Image: "docker.io/ruby"},
  Runtime{Name: "node", Exec: "node", Image: "docker.io/node"}}
