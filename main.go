package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"os"
)

func main() {

	/*



		endpoint := "http://104.131.90.110:4243"
		client, _ := docker.NewClient(endpoint)
		imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
		for _, img := range imgs {
			fmt.Println("ID: ", img.ID)
			fmt.Println("RepoTags: ", img.RepoTags)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Size: ", img.Size)
			fmt.Println("VirtualSize: ", img.VirtualSize)
			fmt.Println("ParentId: ", img.ParentID)
		}*/

	app := cli.NewApp()
	app.Name = "attach"
	app.Usage = "attach docker stdout"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "docker",
			Value: "api",
			Usage: "docker image name to connect",
		},
	}

	app.Action = func(c *cli.Context) {

		endpoint := "http://104.131.90.110:4243"
		client, _ := docker.NewClient(endpoint)
		imgs, _ := client.ListContainers(docker.ListContainersOptions{All: false})

		for _, img := range imgs {
			fmt.Println("ID: ", img.ID)
			fmt.Println("Image: ", img.Image)
			fmt.Println("Command: ", img.Command)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Status: ", img.Status)
			fmt.Println("Names: ", img.Names)
		}

		//r, w := io.Pipe()

		err := client.AttachToContainer(docker.AttachToContainerOptions{Container: "c1c1d16b2353113b6ce73ef41442703f952a6efa0417a102ea4e21112df22332", OutputStream: os.Stdout, InputStream: os.Stdin, Logs: true, Stream: true, Stdin: true, Stdout: true, Stderr: true, RawTerminal: true})

		if err != nil {

			fmt.Printf("%s", err)

		}

	}

	app.Run(os.Args)

}
