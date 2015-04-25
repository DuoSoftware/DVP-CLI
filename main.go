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

	app.Commands = []cli.Command{
		{

			Name:  "attach",
			Usage: "attach docker stdout",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "docker",
					Value: "api",
					Usage: "docker image name to connect",
				},
				cli.StringFlag{
					Name:  "protocol",
					Value: "http",
					Usage: "docker remote api protocol to connect",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1",
					Usage: "docker host ip",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "4243",
					Usage: "docker host port",
				},
				cli.StringFlag{
					Name:  "unixsocket",
					Value: "var/run/docker.sock",
					Usage: "docker unix socket path",
				},
			},

			Action: func(c *cli.Context) {

				name := c.String("docker")
				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				client, _ := docker.NewClient(endpoint)
				imgs, _ := client.ListContainers(docker.ListContainersOptions{All: false})

				imageID := "0"
				imageFound := false

				for _, img := range imgs {
					fmt.Println("ID: ", img.ID)
					fmt.Println("Image: ", img.Image)
					fmt.Println("Command: ", img.Command)
					fmt.Println("Created: ", img.Created)
					fmt.Println("Status: ", img.Status)
					fmt.Println("Names: ", img.Names)

					for _, a := range img.Names {
						if a == name {
							imageFound = true
							fmt.Printf("Image found ....... %s", img.ID)
							imageID = img.ID
							break
						}
					}

					if imageFound == true {
						break
					}
				}

				//r, w := io.Pipe()

				if imageFound {
					err := client.AttachToContainer(docker.AttachToContainerOptions{Container: imageID, OutputStream: os.Stdout, InputStream: os.Stdin, Logs: true, Stream: false, Stdin: true, Stdout: true, Stderr: true, RawTerminal: true})
					if err != nil {

						fmt.Printf("%s", err)

					}
				} else {

					fmt.Printf("Image not found ----------------> %s", name)

				}

			},
		},

		{

			Name:  "list",
			Usage: "list dockers in given host",

			Flags: []cli.Flag{

				cli.StringFlag{
					Name:  "protocol",
					Value: "http",
					Usage: "docker remote api protocol to connect",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1",
					Usage: "docker host ip",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "4243",
					Usage: "docker host port",
				},
				cli.StringFlag{
					Name:  "unixsocket",
					Value: "var/run/docker.sock",
					Usage: "docker unix socket path",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

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

			},
		},

		{

			Name:  "log",
			Usage: "get docker logs through stdout and stderr",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "docker",
					Value: "api",
					Usage: "docker image name to connect",
				},
				cli.StringFlag{
					Name:  "protocol",
					Value: "http",
					Usage: "docker remote api protocol to connect",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1",
					Usage: "docker host ip",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "4243",
					Usage: "docker host port",
				},
				cli.StringFlag{
					Name:  "unixsocket",
					Value: "var/run/docker.sock",
					Usage: "docker unix socket path",
				},
			},

			Action: func(c *cli.Context) {

				name := c.String("docker")
				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				client, _ := docker.NewClient(endpoint)
				imgs, _ := client.ListContainers(docker.ListContainersOptions{All: false})

				imageID := "0"
				imageFound := false

				for _, img := range imgs {
					fmt.Println("ID: ", img.ID)
					fmt.Println("Image: ", img.Image)
					fmt.Println("Command: ", img.Command)
					fmt.Println("Created: ", img.Created)
					fmt.Println("Status: ", img.Status)
					fmt.Println("Names: ", img.Names)

					for _, a := range img.Names {
						if a == name {
							imageFound = true
							fmt.Printf("Image found ....... %s", img.ID)
							imageID = img.ID
							break
						}
					}

					if imageFound == true {
						break
					}
				}

				//r, w := io.Pipe()

				if imageFound {
					err := client.Logs(docker.LogsOptions{Container: imageID, OutputStream: os.Stdout, Stdout: true, Stderr: true, Timestamps: true})
					if err != nil {

						fmt.Printf("%s", err)

					}
				} else {

					fmt.Printf("Image not found ----------------> %s", name)

				}

			},
		},
	}

	app.Run(os.Args)

}
