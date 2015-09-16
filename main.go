package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"github.com/jmcvetta/restclient"
	"os"
	"sort"
	"strings"
)

///////////////////////////swarm cluster /////////////////////////////////////////

type SwarmInstanceOut struct {
	Name      string
	ParentApp string
	UUID      string
	Code      int
	Company   int
	Tenant    int
	Class     string
	Type      string
	Category  string
}

type SwarmInstanceIn struct {
	Name      string
	ParentApp string
	UUID      string
	Code      int
	Company   int
	Tenant    int
	Class     string
	Type      string
	Category  string
	NodeName  string
}

type SwarmNodeOut struct {
	Name           string
	Status         bool
	Code           string
	Company        int
	Tenant         int
	Class          string
	Type           string
	Category       string
	MainIP         string
	Domain         string
	HostDomain     string
	SwarmInstances []SwarmInstanceOut
}

type SwarmNodeIn struct {
	Name         string
	Status       bool
	Code         string
	Company      int
	Tenant       int
	Class        string
	Type         string
	Category     string
	MainIP       string
	Domain       string
	HostDomain   string
	ClusterToken string
}

type SwarmClusterOut struct {
	Name       string
	Token      string
	Code       string
	Company    int
	Tenant     int
	Class      string
	Type       string
	Category   string
	LBDomain   string
	LBIP       string
	SwarmNodes []SwarmNodeOut
}

type SwarmClusterIn struct {
	Name     string
	Token    string
	Code     string
	Company  int
	Tenant   int
	Class    string
	Type     string
	Category string
	LBDomain string
}

//////////////////////////////////////////////////////////////////////////////////

type Service struct {
	Name             string
	Description      string
	Class            string
	Type             string
	Category         string
	CompanyId        int
	TenantId         int
	MultiPorts       bool
	Direction        string
	Protocol         string
	DefaultStartPort int
}

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
	Services           []Service
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

type ENV struct {
	Name   string
	Value  string
	Export bool
}

type Port struct {
	Name  string
	Value string
	Link  bool
}

type Instance struct {
	Name     string
	Class    string
	Type     string
	Category string
	UUID     string
	Image    string
	Ports    []Port
	Envs     []ENV
}

type Deployment struct {
	Name           string
	InternalDomain string
	Class          string
	Type           string
	Category       string
	CompanyId      int
	TenantId       int
	Template       string
	PublicIP       string
	PublicDomain   string
	UUID           string
	Instances      []Instance
}

type Images []Image

func (s Template) Len() int {
	return len(s.TemplateImage)
}
func (s Template) Swap(i, j int) {
	s.TemplateImage[i], s.TemplateImage[j] = s.TemplateImage[j], s.TemplateImage[i]
}
func (s Template) Less(i, j int) bool {
	return s.TemplateImage[i].CSDB_TemplateImage.Priority < s.TemplateImage[j].CSDB_TemplateImage.Priority
}

func main() {

	/*
		timex := time.Date(2015, 7, 24, 12, 53, 0, 0, time.UTC)
		fmt.Printf(timex.Local().String())
		fmt.Printf(time.Now().String())
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

			Name:  "install-instance",
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
					Name:  "sysregistryhost",
					Value: "127.0.0.1",
					Usage: "registry ip",
				},
				cli.StringFlag{
					Name:  "sysregistryport",
					Value: "4243",
					Usage: "registry port",
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
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
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

				url := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/TemplateByName/%s", reghost, regport, template)

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

								fmt.Printf("Template found ready to install --> Enter Deployment name\n")

								reader := bufio.NewReader(os.Stdin)
								text, _ := reader.ReadString('\n')
								fmt.Println(text)
								t := strings.TrimSpace(text)

								fmt.Printf("Template found ready to install --> Enter internal domain name\n")

								texxt, _ := reader.ReadString('\n')
								fmt.Println(texxt)
								d := strings.TrimSpace(texxt)

								fmt.Printf("Enter public domain name\n")

								texyt, _ := reader.ReadString('\n')
								fmt.Println(texyt)
								p := strings.TrimSpace(texyt)

								dep := Deployment{Name: t, InternalDomain: d}
								dep.Class = "USER"
								dep.Type = "DOCKER"
								dep.Category = "SINLEHOST"
								dep.Template = temp.Name
								dep.InternalDomain = d
								dep.PublicIP = host
								dep.PublicDomain = fmt.Sprintf("%s.xip.io", host)
								dep.PublicDomain = p

								sort.Sort(temp)

								for _, img := range temp.TemplateImage {
									fmt.Println(img.CSDB_TemplateImage.Type)
									isInstall := false

									if img.CSDB_TemplateImage.Type == "Mandetory" {
										isInstall = true

									} else if img.CSDB_TemplateImage.Type == "Optional" {

										fmt.Printf("%s %s", img.Name, img.Description)
										fmt.Printf("Above Service is optional do you want to install it?")

										reader := bufio.NewReader(os.Stdin)
										text, _ := reader.ReadString('\n')
										fmt.Println(text)
										t := strings.TrimSpace(text)

										if t == "y" {

											isInstall = true
											fmt.Printf("Install is true")

										} else {
											fmt.Printf("Install is false")
										}

									}

									if isInstall {
										fmt.Println(img.Class)

										ins := Instance{Name: img.Name}
										ins.Class = img.Class
										ins.Type = img.Type
										ins.Category = img.Category

										if img.Class == "DOCKER" {

											fmt.Println(img.DockerUrl)

											pullImage := docker.PullImageOptions{Repository: img.DockerUrl, Tag: "latest"}
											authConf := docker.AuthConfiguration{}
											erry := client.PullImage(pullImage, authConf)
											fmt.Printf("pull --->", erry)
											if erry == nil {

											}

										} else if img.Class == "DOCKERFILE" {

											buildOption := docker.BuildImageOptions{Name: img.Name, Dockerfile: "Dockerfile", SuppressOutput: true, OutputStream: os.Stdout, Remote: img.SourceUrl}

											erry := client.BuildImage(buildOption)
											fmt.Printf("BuildContainer --->", err)

											if erry == nil {

											}

										}

										container := docker.CreateContainerOptions{}
										container.Name = img.Name
										//img.Cmd = "postgres"
										//cmd := []string{img.Cmd}

										/*
											a := []int{1,2,3}
											a = append(a, 4)
											fmt.Println(a)

										*/

										Var := []string{}

										Var = append(Var, fmt.Sprintf("DEPLOYMENT_ENV=%s", "docker"))
										Var = append(Var, fmt.Sprintf("HOST_NAME=%s", img.Name))
										Var = append(Var, fmt.Sprintf("HOST_VERSION=%s", img.Version))

										fmt.Printf("..........................\n", img.SystemVariables)

										for _, vars := range img.SystemVariables {

											envx := ENV{}
											envx.Name = vars.Name
											envx.Export = vars.Export

											fmt.Printf("------------>\n", vars.Type)

											varValue := vars.DefaultValue

											if vars.Type == "uservariable" {

												fmt.Printf("Please enter value for ENV %s ", vars.Name)
												reader := bufio.NewReader(os.Stdin)
												text, _ := reader.ReadString('\n')
												fmt.Println(text)

												enterValue := strings.TrimSpace(text)
												if len(enterValue) > 0 {

													varValue = enterValue

												}

											}

											envx.Value = varValue
											ins.Envs = append(ins.Envs, envx)

											Var = append(Var, fmt.Sprintf("%s=%s", vars.Name, varValue))
										}

										ports := make(map[docker.Port]struct{})

										/////////////////////////////Service Management////////////////////
										for _, servs := range img.Services {

											por := Port{}
											if servs.Direction == "OUT" {
												por.Link = true
											} else {
												por.Link = false
											}

											por.Name = fmt.Sprintf("SYS_%s_%s", servs.Category, servs.Type)

											por.Value = fmt.Sprintf("%d", servs.DefaultStartPort)

											ins.Ports = append(ins.Ports, por)

											Var = append(Var, fmt.Sprintf("HOST_%s_%s=%d", servs.Category, servs.Type, servs.DefaultStartPort))

											portVar := docker.Port(fmt.Sprintf("%d/%s", servs.DefaultStartPort, servs.Protocol))

											var se struct{}
											ports[portVar] = se

										}

										///////////////////////////////////////////////////////////////////////

										/////////////////////////////////Dependancy Management////////////////////

										for _, depe := range img.Dependants {

											itemFound := false

											for _, serchint := range dep.Instances {

												if depe.Name == serchint.Name {

													itemFound = true
													Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s.%s", depe.Category, "HOST", depe.Name, dep.InternalDomain))

													for _, envx := range serchint.Envs {

														Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s", depe.Category, envx.Name, envx.Value))
													}

													for _, portx := range serchint.Ports {

														Var = append(Var, fmt.Sprintf("%s=%s", portx.Name, portx.Value))
													}

													break
												}

											}

											if !itemFound {

												fmt.Println("Dependency Instance %s NotFound ----------->", depe.Name)

											}

										}

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//Cmd: cmd,
										Var = append(Var, fmt.Sprintf("VIRTUAL_HOST=%s.%s", img.Name, dep.PublicDomain))

										fmt.Println("All VARS ----------->", Var)
										container.Config = &docker.Config{Image: img.Name, Env: Var}

										fmt.Println(container.Config.Image)

										containerInstance, errx := client.CreateContainer(container)

										fmt.Printf("Container --->", errx, containerInstance)

										if errx == nil {

											//fmt.Printf("Container ---> ", cont)

											hostConfig := &docker.HostConfig{}

											errz := client.StartContainer(img.Name, hostConfig)
											fmt.Printf("Container --->", errz)
										}

										////////////////////////Add Instance/////////////////////////////////////////////

										dep.Instances = append(dep.Instances, ins)

										/////////////////////////////////////////////////////////////////////

									}

									b, err := json.Marshal(dep)
									if err != nil {
										fmt.Println(err)
										return
									}

									f, err := os.Create("tempfile")

									f.Write(b)
									f.Close()
								}
							}
						}
					}
				}
			},
		},
		///////////////////////////////////////////install-cluster/////////////////////////////////////////////
		{

			Name:  "install-cluster",
			Usage: "install templates in given cluster",

			Flags: []cli.Flag{

				cli.StringFlag{
					Name:  "protocol",
					Value: "http",
					Usage: "docker remote api protocol to connect",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1",
					Usage: "docker swarm master ip",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "4243",
					Usage: "docker swarm master port",
				},
				cli.StringFlag{
					Name:  "token",
					Value: "",
					Usage: "docker cluster token",
				},
				cli.StringFlag{
					Name:  "sysregistryhost",
					Value: "127.0.0.1",
					Usage: "registry ip",
				},
				cli.StringFlag{
					Name:  "sysregistryport",
					Value: "4243",
					Usage: "registry port",
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
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
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

				url := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/TemplateByName/%s", reghost, regport, template)

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

								fmt.Printf("Template found ready to install --> Enter Deployment name\n")

								reader := bufio.NewReader(os.Stdin)
								text, _ := reader.ReadString('\n')
								fmt.Println(text)
								t := strings.TrimSpace(text)

								fmt.Printf("Template found ready to install --> Enter internal domain name\n")

								texxt, _ := reader.ReadString('\n')
								fmt.Println(texxt)
								d := strings.TrimSpace(texxt)

								fmt.Printf("Enter public domain name\n")

								texyt, _ := reader.ReadString('\n')
								fmt.Println(texyt)
								p := strings.TrimSpace(texyt)

								dep := Deployment{Name: t, InternalDomain: d}
								dep.Class = "USER"
								dep.Type = "DOCKER"
								dep.Category = "SINLEHOST"
								dep.Template = temp.Name
								dep.InternalDomain = d
								dep.PublicIP = host
								dep.PublicDomain = fmt.Sprintf("%s.xip.io", host)
								dep.PublicDomain = p

								sort.Sort(temp)

								for _, img := range temp.TemplateImage {
									fmt.Println(img.CSDB_TemplateImage.Type)
									isInstall := false

									if img.CSDB_TemplateImage.Type == "Mandetory" {
										isInstall = true

									} else if img.CSDB_TemplateImage.Type == "Optional" {

										fmt.Printf("%s %s", img.Name, img.Description)
										fmt.Printf("Above Service is optional do you want to install it?")

										reader := bufio.NewReader(os.Stdin)
										text, _ := reader.ReadString('\n')
										fmt.Println(text)
										t := strings.TrimSpace(text)

										if t == "y" {

											isInstall = true
											fmt.Printf("Install is true")

										} else {
											fmt.Printf("Install is false")
										}

									}

									if isInstall {
										fmt.Println(img.Class)

										ins := Instance{Name: img.Name}
										ins.Class = img.Class
										ins.Type = img.Type
										ins.Category = img.Category

										if img.Class == "DOCKER" {

											fmt.Println(img.DockerUrl)

											pullImage := docker.PullImageOptions{Repository: img.DockerUrl, Tag: "latest"}
											authConf := docker.AuthConfiguration{}
											erry := client.PullImage(pullImage, authConf)
											fmt.Printf("pull --->", erry)
											if erry == nil {

											}

										} else if img.Class == "DOCKERFILE" {

											buildOption := docker.BuildImageOptions{Name: img.Name, Dockerfile: "Dockerfile", SuppressOutput: true, OutputStream: os.Stdout, Remote: img.SourceUrl}

											erry := client.BuildImage(buildOption)
											fmt.Printf("BuildContainer --->", err)

											if erry == nil {

											}

										}

										container := docker.CreateContainerOptions{}
										container.Name = img.Name
										//img.Cmd = "postgres"
										//cmd := []string{img.Cmd}

										/*
											a := []int{1,2,3}
											a = append(a, 4)
											fmt.Println(a)

										*/

										Var := []string{}

										Var = append(Var, fmt.Sprintf("DEPLOYMENT_ENV=%s", "docker"))
										Var = append(Var, fmt.Sprintf("HOST_NAME=%s", img.Name))
										Var = append(Var, fmt.Sprintf("HOST_VERSION=%s", img.Version))

										fmt.Printf("..........................\n", img.SystemVariables)

										for _, vars := range img.SystemVariables {

											envx := ENV{}
											envx.Name = vars.Name
											envx.Export = vars.Export

											fmt.Printf("------------>\n", vars.Type)

											varValue := vars.DefaultValue

											if vars.Type == "uservariable" {

												fmt.Printf("Please enter value for ENV %s ", vars.Name)
												reader := bufio.NewReader(os.Stdin)
												text, _ := reader.ReadString('\n')
												fmt.Println(text)

												enterValue := strings.TrimSpace(text)
												if len(enterValue) > 0 {

													varValue = enterValue

												}

											}

											envx.Value = varValue
											ins.Envs = append(ins.Envs, envx)

											Var = append(Var, fmt.Sprintf("%s=%s", vars.Name, varValue))
										}

										ports := make(map[docker.Port]struct{})

										/////////////////////////////Service Management////////////////////
										for _, servs := range img.Services {

											por := Port{}
											if servs.Direction == "OUT" {
												por.Link = true
											} else {
												por.Link = false
											}

											por.Name = fmt.Sprintf("SYS_%s_%s", servs.Category, servs.Type)

											por.Value = fmt.Sprintf("%d", servs.DefaultStartPort)

											ins.Ports = append(ins.Ports, por)

											Var = append(Var, fmt.Sprintf("HOST_%s_%s=%d", servs.Category, servs.Type, servs.DefaultStartPort))

											portVar := docker.Port(fmt.Sprintf("%d/%s", servs.DefaultStartPort, servs.Protocol))

											var se struct{}
											ports[portVar] = se

										}

										///////////////////////////////////////////////////////////////////////

										/////////////////////////////////Dependancy Management////////////////////

										for _, depe := range img.Dependants {

											itemFound := false

											for _, serchint := range dep.Instances {

												if depe.Name == serchint.Name {

													itemFound = true
													Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s.%s", depe.Category, "HOST", depe.Name, dep.InternalDomain))

													for _, envx := range serchint.Envs {

														Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s", depe.Category, envx.Name, envx.Value))
													}

													for _, portx := range serchint.Ports {

														Var = append(Var, fmt.Sprintf("%s=%s", portx.Name, portx.Value))
													}

													break
												}

											}

											if !itemFound {

												fmt.Println("Dependency Instance %s NotFound ----------->", depe.Name)

											}

										}

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//Cmd: cmd,
										Var = append(Var, fmt.Sprintf("VIRTUAL_HOST=%s.%s", img.Name, dep.PublicDomain))

										fmt.Println("All VARS ----------->", Var)
										container.Config = &docker.Config{Image: img.Name, Env: Var}

										fmt.Println(container.Config.Image)

										containerInstance, errx := client.CreateContainer(container)

										fmt.Printf("Container --->", errx, containerInstance)

										if errx == nil {

											//fmt.Printf("Container ---> ", cont)

											hostConfig := &docker.HostConfig{}

											errz := client.StartContainer(img.Name, hostConfig)
											fmt.Printf("Container --->", errz)
										}

										////////////////////////Add Instance/////////////////////////////////////////////

										dep.Instances = append(dep.Instances, ins)

										/////////////////////////////////////////////////////////////////////

									}

									b, err := json.Marshal(dep)
									if err != nil {
										fmt.Println(err)
										return
									}

									f, err := os.Create("tempfile")

									f.Write(b)
									f.Close()
								}
							}
						}
					}
				}
			},
		},

		/////////////////////////////////////////////////////////////////////////////////////////

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
