package main

import (
	//"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"github.com/jmcvetta/restclient"
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

			Name:  "install",
			Usage: "install templates in given host",

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
				cli.StringFlag{
					Name:  "template",
					Value: "DBTEMPLATE",
					Usage: "template name to be install",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				template := c.String("template")

				fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				client, _ := docker.NewClient(endpoint)
				//client.ListContainers(docker.ListContainersOptions{All: false})

				/////////--------------------------------------------------------------------------------------------->

				type Variable struct {
					Name         string
					Description  string
					DefaultValue string
					Export       bool
					Type         string
				}

				type Template_Image struct {
					Type     string
					Priority int
				}

				type Image struct {
					Name               string
					Description        string
					Version            string
					VersionStatus      string
					SourceUrl          string
					DockerUrl          string
					Class              string
					Type               string
					Cmd                string
					Category           string
					Importance         string
					CSDB_TemplateImage Template_Image
					SystemVariables    []Variable
					Dependants         []Image
				}

				type Template struct {
					Name          string
					Description   string
					Class         string
					Type          string
					Category      string
					CompanyId     int
					TenantId      int
					TemplateImage []Image
				}

				type Result struct {
					Exception     string
					CustomMessage string
					IsSuccess     bool
					Result        []Template
				}

				url := fmt.Sprintf("http://127.0.0.1:9093/DVP/API/1.0/SystemRegistry/TemplateByName/%s", template)

				var s Result

				r := restclient.RequestResponse{
					Url:    url,
					Method: "GET",
					Result: &s,
				}
				status, err := restclient.Do(&r)
				if err != nil {
					//panic(err)
				}
				if status == 200 {

					//json.Unmarshal([]byte(r.RawText), &s)

					fmt.Println("Template Data  -->", r.RawText)

					fmt.Println("Template Data  -->", s)

					if s.IsSuccess == true {

						if s.Result != nil {

							for _, temp := range s.Result {

								for _, img := range temp.TemplateImage {
									fmt.Println(img.CSDB_TemplateImage.Type)
									if img.CSDB_TemplateImage.Type == "Mandetory" {

										fmt.Println(img.Class)
										if img.Class == "DOCKER" {

											fmt.Println(img.DockerUrl)

											pullImage := docker.PullImageOptions{Repository: img.DockerUrl, Tag: "latest"}
											authConf := docker.AuthConfiguration{}
											erry := client.PullImage(pullImage, authConf)
											fmt.Printf("pull --->", erry)
											if erry == nil {

												container := docker.CreateContainerOptions{}
												container.Name = img.Name
												//img.Cmd = "postgres"
												cmd := []string{img.Cmd}

												/*
													a := []int{1,2,3}
													a = append(a, 4)
													fmt.Println(a)

												*/

												Var := []string{}

												for _, vars := range img.SystemVariables {

													Var = append(Var, fmt.Sprintf("%s=%s", vars.Name, vars.DefaultValue))

												}

												container.Config = &docker.Config{Image: img.DockerUrl, Cmd: cmd, Env: Var}

												fmt.Println(container.Config.Image)

												_, errx := client.CreateContainer(container)

												fmt.Printf("Container --->", errx)

												if errx == nil {

													//fmt.Printf("Container ---> ", cont)

													hostConfig := &docker.HostConfig{}

													errz := client.StartContainer(img.Name, hostConfig)
													fmt.Printf("Container --->", errz)
												}
											}

										} else if img.Class == "DOCKERFILE" {

											buildOption := docker.BuildImageOptions{Name: img.Name, Dockerfile: "Dockerfile", SuppressOutput: true, OutputStream: os.Stdout, Remote: img.SourceUrl}

											err := client.BuildImage(buildOption)
											fmt.Printf("BuildContainer --->", err)

										}

									} else if img.CSDB_TemplateImage.Type == "Optional" {

									}

								}

							}

						}
					}

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
