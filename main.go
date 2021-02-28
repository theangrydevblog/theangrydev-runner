package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	containerStruct "./container"
	"./runtime"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var ctx context.Context = context.Background()
var wg sync.WaitGroup

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	runtimes := runtime.Runtimes
	c := make(chan os.Signal)
	alertChannel := make(chan containerStruct.Health)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nSIGNINT handler invoked. Forecfully removing all containers")
		for _, r := range runtimes {
			r.Container.RemoveContainer(ctx, cli)
		}

	}()

	wg.Add(len(runtimes))

	for _, r := range runtimes {
		go func(r *runtime.Runtime) {
			if r.CheckIfExistsLocally(ctx, cli) != true {
				fmt.Printf("%s not found in host. Pulling from container registry...\n", r.Tag)
				out, pullError := cli.ImagePull(ctx, r.Tag, types.ImagePullOptions{All: false})
				if pullError != nil {
					panic(pullError)
				}

				defer out.Close()
				io.Copy(os.Stdout, out)
			}

			if r.Container.CheckIfRunning(ctx, cli) {
				fmt.Printf("%s is already running. Grabbing ID...\n", r.Container.Name)
				r.Container.UpdateID(ctx, cli)
			} else {
				r.Container.StartContainer(ctx, cli, r.Image)
			}

			wg.Done()
		}(r)
	}

	wg.Wait()

	pyData, err := ioutil.ReadFile("./tests/scripts/python.py")
	if err != nil {
		panic(err)
	}

	python := runtimes["Python"]
	pyOutput := python.ExecuteCode(ctx, cli, string(pyData))
	fmt.Println(pyOutput)

	rbData, err := ioutil.ReadFile("./tests/scripts/ruby.rb")
	if err != nil {
		panic(err)
	}

	ruby := runtimes["Ruby"]
	rbOutput := ruby.ExecuteCode(ctx, cli, string(rbData))
	fmt.Println(rbOutput)

	for _, r := range runtimes {
		fmt.Printf("Firing off health checker for %s\n", r.Container.Name)
		go r.Container.Monitor(ctx, cli, alertChannel)
	}

	fmt.Println("Waiting for requests...")
	for {
		message := <-alertChannel
		if message.MessageCode != true {
			fmt.Printf("Container %s is shut down...restarting\n", message.Container.Name)
			go message.Container.RestartContainer(ctx, cli)
		}

	}

}
