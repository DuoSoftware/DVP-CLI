package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"github.com/jmcvetta/restclient"
	"github.com/satori/go.uuid"
	"os"
	"sort"
	"strings"
)

///////////////////////////swarm cluster /////////////////////////////////////////

type SwarmInstanceOut struct {
	Name        string
	ParentApp   string
	UUID        string
	Code        string
	FrontEnd    string
	BackEnd     string
	SwarmNodeId string
	Company     int
	Tenant      int
	Class       string
	Type        string
	Category    string
}

type SwarmInstanceIn struct {
	SwarmNodeUuid  string
	DeploymentUUID string
	Name           string
	ParentApp      string
	UUID           string
	Code           string
	FrontEnd       string
	BackEnd        string
	Company        int
	Tenant         int
	Class          string
	Type           string
	Category       string
	NodeName       string
	Envs           []ENV
	Ports          []Port
}

type SwarmNodeOut struct {
	UUID          string
	Name          string
	Status        bool
	Code          string
	Company       int
	Tenant        int
	Class         string
	Type          string
	Category      string
	MainIP        string
	Domain        string
	HostDomain    string
	SwarmInstance []SwarmInstanceOut
}

type SwarmNodeOutx struct {
	UUID       string
	Name       string
	Status     bool
	Code       string
	Company    int
	Tenant     int
	Class      string
	Type       string
	Category   string
	MainIP     string
	Domain     string
	HostDomain string
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
	Name      string
	Token     string
	Code      int
	Company   int
	Tenant    int
	Class     string
	Type      string
	Category  string
	LBDomain  string
	LBIP      string
	SwarmNode []SwarmNodeOut
}

type SwarmClusterIn struct {
	Name     string
	Token    string
	Code     int
	Company  int
	Tenant   int
	Class    string
	Type     string
	Category string
	LBDomain string
}

type ClusterResult struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        SwarmClusterOut
}

type BasicResult struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
}

type InstanceResult struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        SwarmInstanceOut
}

type NodeResult struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        []SwarmNodeOutx
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

type ImageResult struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        Image
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
	Host     string
	Name     string
	Class    string
	Type     string
	Category string
	UUID     string
	Image    string
	State    string
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

//type Images []Image

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

		//////////////////////////////////////////////attach/////////////////////////////////////////////////////////
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

		////////////////////////////////////////////install-instance/////////////////////////////////////////////////
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
								fmt.Println("temp.TemplateImage", len(temp.TemplateImage), "%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
								for _, img := range temp.TemplateImage {
									fmt.Println(img.CSDB_TemplateImage.Type)
									isInstall := false

									if img.CSDB_TemplateImage.Type == "Mandetory" {
										isInstall = true
										fmt.Println("IsInstall:", isInstall)

									} else if img.CSDB_TemplateImage.Type == "Optional" {

										fmt.Println("%s %s", img.Name, img.Description)
										fmt.Println("Above Service is optional do you want to install it?")

										reader := bufio.NewReader(os.Stdin)
										text, _ := reader.ReadString('\n')
										fmt.Println(text)
										t := strings.TrimSpace(text)

										if t == "y" {

											isInstall = true
											fmt.Println("Install is true")

										} else {
											fmt.Println("Install is false")
										}

									}

									if isInstall {
										fmt.Println("Start Install: ", img.Class)

										ins := Instance{Name: img.Name}
										ins.Class = img.Class
										ins.Type = img.Type
										ins.Category = img.Category

										if img.Class == "DOCKER" {

											fmt.Println(img.DockerUrl)

											pullImage := docker.PullImageOptions{Repository: img.DockerUrl, Tag: "latest"}
											authConf := docker.AuthConfiguration{}
											erry := client.PullImage(pullImage, authConf)
											fmt.Println("pull --->", erry)
											if erry == nil {

											}

										} else if img.Class == "DOCKERFILE" {

											buildOption := docker.BuildImageOptions{Name: img.Name, Dockerfile: "Dockerfile", SuppressOutput: true, OutputStream: os.Stdout, Remote: img.SourceUrl}

											erry := client.BuildImage(buildOption)
											fmt.Println("BuildContainer --->", err)

											if erry == nil {

											}

										}
										fmt.Println("Start docker.CreateContainerOptions %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
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

										fmt.Println("..........................\n", img.SystemVariables)

										for _, vars := range img.SystemVariables {

											envx := ENV{}
											envx.Name = vars.Name
											envx.Export = vars.Export

											//fmt.Printf("------------>\n", vars.Type)

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

										fmt.Println("Start Service Management %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
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
										fmt.Println("Start Dependancy Management %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
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

												fmt.Printf("Dependency Instance %s NotFound -----------> \n", depe.Name)

											}

										}
										fmt.Println("End Session Management %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//Cmd: cmd,
										Var = append(Var, fmt.Sprintf("VIRTUAL_HOST=%s.%s", img.Name, dep.PublicDomain))

										fmt.Println("All VARS ----------->", Var)
										container.Config = &docker.Config{Image: img.Name, Env: Var}

										fmt.Println(container.Config.Image)

										containerInstance, errx := client.CreateContainer(container)

										fmt.Println("Container --->", errx, containerInstance)

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

		////////////////////////////////////////////update-image////////////////////////////////////////////////////
		{
			Name:  "update-image",
			Usage: "update image in swarn cluster",

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
					Name:  "lbapihost",
					Value: "127.0.0.1",
					Usage: "Hipache API host",
				},

				cli.StringFlag{
					Name:  "containername",
					Value: "",
					Usage: "Container Name",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
				id := c.String("containername")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				if len(id) > 0 {

					//http://45.55.142.207:8826/DVP/API/1.0.0.0/SystemRegistry/ImageByName/fileservice
					urlx := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/ImageByName/%s", reghost, regport, id)

					var imgres ImageResult

					rx := restclient.RequestResponse{
						Url:    urlx,
						Method: "GET",
						Result: &imgres,
					}

					statusx, errx := restclient.Do(&rx)

					img := imgres.Result

					fmt.Printf("%s -> %d %s", urlx, statusx, img.Name)

					if errx != nil {
						//panic(err)
					}
					if statusx == 200 {

						client, _ := docker.NewClient(endpoint)

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
							fmt.Printf("BuildContainer --->", erry)

							if erry == nil {

							}

						}

					}

				} else {

					fmt.Printf("Container id is required")

				}

			},
		},

		///////////////////////////////////////////install-cluster///////////////////////////////////////////////////
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

				cli.BoolFlag{
					Name:  "defaultvar",
					Usage: "run with default values",
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
				dockerClusterToken := c.String("token")
				defaultvar := c.Bool("defaultvar")

				fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				fmt.Println("Endpoint ------------------------------------> ", endpoint)
				client, _ := docker.NewClient(endpoint)
				//client.ListContainers(docker.ListContainersOptions{All: false})
				/////////--------------------------------------------------------------------------------------------->
				url := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/TemplateByName/%s", reghost, regport, template)
				cUrl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/ClusterByToken/%s", reghost, regport, dockerClusterToken)

				var s Result
				var cs ClusterResult

				r := restclient.RequestResponse{
					Url:    url,
					Method: "GET",
					Result: &s,
				}

				cr := restclient.RequestResponse{
					Url:    cUrl,
					Method: "GET",
					Result: &cs,
				}
				status, err := restclient.Do(&r)
				cStatus, err := restclient.Do(&cr)
				if err != nil {
					//panic(err)
				}
				if status == 200 && cStatus == 200 {

					//json.Unmarshal([]byte(r.RawText), &s)

					fmt.Println("Template Data  -->", r.RawText)

					fmt.Println("Template Data  -->", s)

					fmt.Println("Cluster Data  -->", cr.RawText)

					fmt.Println("Cluster Data  -->", cs)

					if s.IsSuccess == true && cs.IsSuccess == true {

						if s.Result != nil {

							for _, temp := range s.Result {

								fmt.Printf("Template found ready to install --> Enter Deployment name\n")

								reader := bufio.NewReader(os.Stdin)
								text, _ := reader.ReadString('\n')
								fmt.Println(text)
								t := strings.TrimSpace(text)

								//fmt.Printf("Template found ready to install --> Enter internal domain name\n")// cluster domain

								//texxt, _ := reader.ReadString('\n')
								//fmt.Println(texxt)
								fmt.Println(cs.Result.LBDomain)
								d := strings.TrimSpace(cs.Result.LBDomain)

								//fmt.Printf("Enter public domain name\n") //

								//texyt, _ := reader.ReadString('\n')
								fmt.Println(cs.Result.LBDomain)
								p := strings.TrimSpace(cs.Result.LBDomain)

								dep := Deployment{Name: t, InternalDomain: d}
								dep.Class = "USER"
								dep.Type = "DOCKER"
								dep.Category = "CLUSTER"
								dep.Template = temp.Name
								dep.InternalDomain = d
								dep.PublicIP = host
								dep.PublicDomain = fmt.Sprintf("%s.xip.io", host)

								u1 := uuid.NewV4()
								dep.UUID = u1.String()

								dep.PublicDomain = p

								sort.Sort(temp)

								fmt.Printf("Images : %v", temp.TemplateImage)

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

									} else if img.CSDB_TemplateImage.Type == "Backing" {

										isInstall = false

										fmt.Println(img.Class)

										ins := Instance{Name: img.Name}
										ins.Class = img.Class
										ins.Type = img.Type
										ins.Category = img.Category

										fmt.Printf("Please enter host name for ", img.Name)
										reader := bufio.NewReader(os.Stdin)
										text, _ := reader.ReadString('\n')
										fmt.Println(text)

										textHost := strings.TrimSpace(text)

										ins.Host = textHost

										///////////////////////////////////////////////
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

										}

										///////////////////////////////////////////////////////////////////////

										//Var = append(Var, fmt.Sprintf("VIRTUAL_HOST=%s.%s", img.Name, dep.PublicDomain))

										dep.Instances = append(dep.Instances, ins)

									}

									if isInstall {
										fmt.Println(img.Class)

										ins := Instance{Name: img.Name}
										idata := SwarmInstanceIn{}

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

												if !defaultvar {
													fmt.Printf("Please enter value for ENV %s ", vars.Name)
													reader := bufio.NewReader(os.Stdin)
													text, _ := reader.ReadString('\n')
													fmt.Println(text)

													enterValue := strings.TrimSpace(text)
													if len(enterValue) > 0 {

														varValue = enterValue

													}
												}

											}

											envx.Value = varValue
											ins.Envs = append(ins.Envs, envx)
											idata.Envs = append(idata.Envs, envx)

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
											idata.Ports = append(idata.Ports, por)

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

													if serchint.Host != "" {

														Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s", depe.Category, "HOST", serchint.Host))
													} else {

														Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s.%s", depe.Category, "HOST", depe.Name, dep.InternalDomain))
													}
													//Var = append(Var, fmt.Sprintf("SYS_%s_%s=%s.%s", depe.Category, "HOST", depe.Name, dep.InternalDomain))

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

												fmt.Printf("Dependency Instance %s NotFound ----------->\n", depe.Name)

											}

										}

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//Cmd: cmd,
										Var = append(Var, fmt.Sprintf("VIRTUAL_HOST=%s.*", img.Name))
										Var = append(Var, fmt.Sprintf("LB_FRONTEND=%s.%s", img.Name, cs.Result.LBDomain))
										Var = append(Var, fmt.Sprintf("LB_PORT=%d", 80))

										fmt.Println("All VARS ----------->", Var)
										container.Config = &docker.Config{Image: img.Name, Env: Var}

										fmt.Println(container.Config.Image)
										fmt.Println("Start Create Container%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

										containerInstanceId, errx := client.CreateContainer(container)

										if errx != nil {

											fmt.Printf("CreateContainer %v \n", errx)
											//break
										} else {

											containerInstance, errx := client.InspectContainer(containerInstanceId.ID)

											if errx != nil {

												fmt.Printf("CreateContainer %v \n", errx)
											}

											fmt.Println("Container ---> %v", containerInstance)
											fmt.Println("End Create Container%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

											fmt.Println("Error not found%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

											//--------------------------------------------------------------------------------------------
											var hdomain string
											fmt.Println("Start find hdomain", len(cs.Result.SwarmNode))
											fmt.Println("%%%%%%%%%%%% SwarmNode %%%%%%%%%%%%%%%%%%%%%")
											fmt.Printf("%v\n", cs.Result.SwarmNode)
											fmt.Println("%%%%%%%%%%%% containerInstance %%%%%%%%%%%%%%%%%%%%%")
											fmt.Printf("%v\n", containerInstance.Node)
											for _, snode := range cs.Result.SwarmNode {
												fmt.Println("ci node id: ", containerInstance.Node.ID)
												fmt.Println("snode.UniqueId: ", snode.UUID)
												if containerInstance.Node.ID == snode.UUID {
													hdomain = snode.Domain
													break
												}
											}
											fmt.Println("End find hdomain", hdomain)
											iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Instance", reghost, regport)
											hUrl := fmt.Sprintf("http://%s:%s/frontends?host=%s.%s&backends=http://%s.%s", cs.Result.LBIP, "5000", img.Name, cs.Result.LBDomain, img.Name, hdomain)

											fmt.Println("Iurl ", iurl, hUrl)

											idata.Class = cs.Result.Class
											idata.Type = cs.Result.Type
											idata.Category = cs.Result.Category
											idata.Company = cs.Result.Company
											idata.Tenant = cs.Result.Tenant
											idata.Code = containerInstance.Name
											idata.NodeName = containerInstance.Node.Name
											idata.ParentApp = img.Name
											idata.SwarmNodeUuid = containerInstance.Node.ID
											idata.DeploymentUUID = dep.UUID
											idata.UUID = containerInstance.ID
											ins.UUID = containerInstance.ID
											ins.State = containerInstance.State.String()
											ins.Image = img.Name
											idata.Name = containerInstance.Name
											idata.FrontEnd = fmt.Sprintf("%s.%s", img.Name, cs.Result.LBDomain)
											idata.BackEnd = fmt.Sprintf("http://%s.%s", img.Name, hdomain)

											var ibs BasicResult
											var hbs string

											ir := restclient.RequestResponse{
												Url:    iurl,
												Method: "POST",
												Data:   &idata,
												Result: &ibs,
											}

											hr := restclient.RequestResponse{
												Url:    hUrl,
												Method: "POST",
												Result: &hbs,
											}
											iStatus, err := restclient.Do(&ir)
											hStatus, err := restclient.Do(&hr)
											if err != nil {
												//panic(err)
											}
											fmt.Println(iStatus, ibs.CustomMessage)
											fmt.Println(hStatus, hbs)
											//------------------------------------------------------------------------------------

											//fmt.Printf("Container ---> ", cont)

											hostConfig := &docker.HostConfig{}
											hostConfig.PublishAllPorts = true

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

								iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/HealthMonitor/Node", reghost, regport)
								var ibs BasicResult

								ir := restclient.RequestResponse{
									Url:    iurl,
									Method: "POST",
									Data:   &dep,
									Result: &ibs,
								}

								iStatus, err := restclient.Do(&ir)

								if err != nil {
									//panic(err)
								}
								fmt.Printf("%s %s %v\n", iurl, iStatus, ibs)

							}

						}

					}
				}
			},
		},

		//////////////////////////////////////////kill-container/////////////////////////////////////////////////////
		{
			Name:  "kill-container",
			Usage: "kill container from swarn cluster",

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
					Name:  "lbapihost",
					Value: "127.0.0.1",
					Usage: "Hipache API host",
				},

				cli.StringFlag{
					Name:  "containerid",
					Value: "",
					Usage: "Container ID",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
				apihost := c.String("lbapihost")
				id := c.String("containerid")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				if len(id) > 0 {

					//http://45.55.142.207:8826/DVP/API/1.0.0.0/SystemRegistry/InstanceById/d52cd7fc09684cf4788b9e3b7cda97f0c9c0d78c0c7fc9a4967bd87de5c0860e

					client, _ := docker.NewClient(endpoint)
					errx := client.KillContainer(docker.KillContainerOptions{ID: id})

					if errx != nil {

						fmt.Printf("Kill container %s is failed", id)

					} else {

						//hUrl := fmt.Sprintf("http://%s:%s/frontends?host=%s.%s&backends=http://%s.%s", cs.Result.LBIP, "5000", img.Name, cs.Result.LBDomain, img.Name, hdomain)
						//DELETE /frontends/:name/backend?backend=http://host1:port

						urlx := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/InstanceById/%s", reghost, regport, id)

						var sx InstanceResult

						rx := restclient.RequestResponse{
							Url:    urlx,
							Method: "GET",
							Result: &sx,
						}

						statusx, errx := restclient.Do(&rx)

						fmt.Printf("%s -> %d", urlx, statusx)

						if errx != nil {
							//panic(err)
						}
						if statusx == 200 {

							////DELETE /frontends/:name/backend?backend=http://host1:port

							hUrl := fmt.Sprintf("http://%s:%s/frontends/%s/backend?backend=%s", apihost, "5000", sx.Result.FrontEnd, sx.Result.BackEnd)

							hr := restclient.RequestResponse{
								Url:    hUrl,
								Method: "DELETE",
							}

							hStatus, erry := restclient.Do(&hr)

							fmt.Printf("%s -> %d", hUrl, hStatus)

							if erry != nil {

							}
							//hStatus == 200
							if hStatus == 200 {

								fmt.Printf("Container successfully deleted %s", id)

								///DVP/API/:version/SystemRegistry/Node/:uuid/Instance/:id/Status/:status

								iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Node/%s/Instance/%s/Status/down", reghost, regport, sx.Result.SwarmNodeId, id)

								ir := restclient.RequestResponse{
									Url:    iurl,
									Method: "PUT",
								}

								iStatus, err := restclient.Do(&ir)

								if err != nil {
									//panic(err)
								}

								fmt.Printf("%s -> %d", iurl, iStatus)
							}
						}
					}

				} else {

					fmt.Printf("Container id is required")

				}

			},
		},

		///////////////////////////////////////////////increase container///////////////////////////////////////////
		{
			Name:  "scale-container",
			Usage: "scale container from swarn cluster",

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
					Name:  "lbapihost",
					Value: "127.0.0.1",
					Usage: "Hipache API host",
				},

				cli.StringFlag{
					Name:  "containerid",
					Value: "",
					Usage: "Container ID",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
				apihost := c.String("lbapihost")
				id := c.String("containerid")

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				if len(id) > 0 {

					urlx := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/InstanceById/%s", reghost, regport, id)

					var sx InstanceResult

					rx := restclient.RequestResponse{
						Url:    urlx,
						Method: "GET",
						Result: &sx,
					}

					statusx, errx := restclient.Do(&rx)

					fmt.Printf("%s -> %d", urlx, statusx)

					fmt.Printf("%v", sx)

					if errx != nil {
						//panic(err)
					}
					if statusx == 200 {

						client, _ := docker.NewClient(endpoint)
						containerx, err := client.InspectContainer(id)

						if err == nil {

							fmt.Printf("container %v %v \n", containerx, err)

							container := docker.CreateContainerOptions{}

							u1 := uuid.NewV4()
							nameuuid := u1.String()

							container.Name = nameuuid

							container.Config = &docker.Config{Image: containerx.Config.Image, Env: containerx.Config.Env}

							cont, errx := client.CreateContainer(container)

							fmt.Printf("container %v %v \n", cont, containerx.Node)

							if errx == nil {

								//cont.Node.ID

								hostConfig := &docker.HostConfig{}
								hostConfig.PublishAllPorts = true

								url := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Node/%s", reghost, regport, containerx.Node.ID)

								var s NodeResult

								r := restclient.RequestResponse{
									Url:    url,
									Method: "GET",
									Result: &s,
								}

								status, err := restclient.Do(&r)

								fmt.Printf("%s -> %d", url, status)

								if err != nil {
									//panic(err)
								}

								if status == 200 {

									errz := client.StartContainer(container.Name, hostConfig)
									fmt.Printf("Container --->", errz)

									fmt.Printf("%v", s.Result)

									hUrl := fmt.Sprintf("http://%s:%s/frontends/%s?backends=http://%s.%s", apihost, "5000", sx.Result.FrontEnd, nameuuid, s.Result[0].Domain)

									var hbs string
									hr := restclient.RequestResponse{
										Url:    hUrl,
										Method: "POST",
										Result: &hbs,
									}

									hStatus, erry := restclient.Do(&hr)

									fmt.Printf("%s -> %d", hUrl, hStatus)

									if erry != nil {

									}

									if hStatus == 200 {

										fmt.Printf("Backend Successfully added")

										idata := SwarmInstanceIn{}
										idata.Class = sx.Result.Class
										idata.Type = sx.Result.Type
										idata.Category = sx.Result.Category
										idata.Company = sx.Result.Company
										idata.Tenant = sx.Result.Tenant
										idata.Code = nameuuid
										idata.NodeName = containerx.Node.Name
										idata.ParentApp = sx.Result.ParentApp
										idata.SwarmNodeUuid = containerx.Node.ID
										idata.DeploymentUUID = cont.ID
										idata.UUID = cont.ID

										idata.Name = cont.Name
										idata.FrontEnd = sx.Result.FrontEnd
										idata.BackEnd = fmt.Sprintf("http://%s.%s", nameuuid, s.Result[0].Domain)

										iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Instance", reghost, regport)

										var ibs BasicResult
										//var hbs string

										ir := restclient.RequestResponse{
											Url:    iurl,
											Method: "POST",
											Data:   &idata,
											Result: &ibs,
										}

										iStatus, err := restclient.Do(&ir)

										if err != nil {

										}

										if iStatus == 200 {

										}
									}
								}
							}
						}
					}
				}
			},
		},

		//////////////////////////////////////////////start-container////////////////////////////////////////////////
		{
			Name:  "start-container",
			Usage: "start container from swarn cluster",

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
					Name:  "lbapihost",
					Value: "127.0.0.1",
					Usage: "Hipache API host",
				},

				cli.StringFlag{
					Name:  "containerid",
					Value: "",
					Usage: "Container ID",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
				apihost := c.String("lbapihost")
				id := c.String("containerid")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				if len(id) > 0 {

					client, _ := docker.NewClient(endpoint)
					errx := client.StartContainer(id, &docker.HostConfig{})
					if errx != nil {

						fmt.Printf("Kill container %s is failed", id)

					} else {

						urlx := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/InstanceById/%s", reghost, regport, id)

						var sx InstanceResult

						rx := restclient.RequestResponse{
							Url:    urlx,
							Method: "GET",
							Result: &sx,
						}

						statusx, errx := restclient.Do(&rx)

						fmt.Printf("%s -> %d", urlx, statusx)

						if errx != nil {
							//panic(err)
						}
						if statusx == 200 {

							//hUrl := fmt.Sprintf("http://%s:%s/frontends?host=%s.%s&backends=http://%s.%s", cs.Result.LBIP, "5000", img.Name, cs.Result.LBDomain, img.Name, hdomain)

							hUrl := fmt.Sprintf("http://%s:%s/frontends?host=%s&backends=%s", apihost, "5000", sx.Result.FrontEnd, sx.Result.BackEnd)

							hr := restclient.RequestResponse{
								Url:    hUrl,
								Method: "POST",
							}

							hStatus, erry := restclient.Do(&hr)

							fmt.Printf("%s -> %d", hUrl, hStatus)

							if erry != nil {

							}
							//hStatus == 200
							if hStatus == 200 {

								fmt.Printf("Container successfully deleted %s", id)

								///DVP/API/:version/SystemRegistry/Node/:uuid/Instance/:id/Status/:status

								iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Node/%s/Instance/%s/Status/up", reghost, regport, sx.Result.SwarmNodeId, id)

								ir := restclient.RequestResponse{
									Url:    iurl,
									Method: "PUT",
								}

								iStatus, err := restclient.Do(&ir)

								if err != nil {
									//panic(err)
								}

								fmt.Printf("%s -> %d", iurl, iStatus)
							}
						}
					}

				} else {

					fmt.Printf("Container id is required")

				}

			},
		},

		/////////////////////////////////////////////////monitor/////////////////////////////////////////////////////
		{
			Name:  "monitor",
			Usage: "monitor helth of swarn cluster",

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
					Name:  "lbapihost",
					Value: "127.0.0.1",
					Usage: "Hipache API host",
				},
			},

			Action: func(c *cli.Context) {

				protocol := c.String("protocol")
				host := c.String("host")
				port := c.String("port")
				socket := c.String("unixsocket")
				reghost := c.String("sysregistryhost")
				regport := c.String("sysregistryport")
				apihost := c.String("lbapihost")

				//fmt.Printf("Image ----------------> %s", c)

				endpoint := fmt.Sprintf("http://%s:%s", host, port)

				if protocol == "unix" {

					endpoint = fmt.Sprintf("unix:///%s", socket)

				}

				client, _ := docker.NewClient(endpoint)
				listener := make(chan *docker.APIEvents)

				client.AddEventListener(listener)

				err := client.AddEventListener(listener)
				if err != nil {
					fmt.Errorf("Failed to add event listener: %s", err)
				}

				//timeout := time.After(1 * time.Second)
				//var count int

				for {
					select {
					case msg := <-listener:

						//busybox node:DUOVOICEAPP2 , stop , a3c37ba2fbfae960825c978bc44b3a6928d0d428fa70fff6a046c97ee228973f , 1443417873
						fmt.Printf("Received: from-%s, Status-%s, ID-%s, Time-%s\n", msg.From, msg.Status, msg.ID, msg.Time)

						///get instance details///////////////////////////////////////////////////////////////////////////////////////////

						url := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/NodeByInstanceId/%s", reghost, regport, msg.ID)

						var s NodeResult

						r := restclient.RequestResponse{
							Url:    url,
							Method: "GET",
							Result: &s,
						}

						status, err := restclient.Do(&r)

						urlx := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/InstanceById/%s", reghost, regport, msg.ID)

						var sx InstanceResult

						rx := restclient.RequestResponse{
							Url:    urlx,
							Method: "GET",
							Result: &sx,
						}

						statusx, errx := restclient.Do(&rx)

						if err != nil || errx != nil {
							//panic(err)
						}
						if status == 200 && statusx == 200 {

							switch msg.Status {

							case "create":

							case "start":

							case "kill":

							case "die":

							case "untag": //image

							case "delete": //image

							case "destroy": //container

								iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Node/%s/Instance/%s", reghost, regport, s.Result[0].UUID, msg.ID)

								var ibs BasicResult

								ir := restclient.RequestResponse{
									Url:    iurl,
									Method: "DELETE",
									Result: &ibs,
								}

								iStatus, err := restclient.Do(&ir)

								if err != nil {
									//panic(err)
								}

								fmt.Println(iurl, iStatus)

							case "stop":

								iurl := fmt.Sprintf("http://%s:%s/DVP/API/1.0/SystemRegistry/Node/%s/Instance/%s/Status/%s", reghost, regport, s.Result[0].UUID, msg.ID, msg.Status)
								hUrl := fmt.Sprintf("http://%s:%s/frontends/%s", apihost, "5000", sx.Result.FrontEnd)

								var ibs BasicResult
								var hbs string

								ir := restclient.RequestResponse{
									Url:    iurl,
									Method: "PUT",
									Result: &ibs,
								}

								hr := restclient.RequestResponse{
									Url:    hUrl,
									Method: "DELETE",
									Result: &hbs,
								}
								iStatus, err := restclient.Do(&ir)
								hStatus, err := restclient.Do(&hr)
								if err != nil {
									//panic(err)
								}

								fmt.Println(iurl, iStatus, hUrl, hStatus)
							default:

							}
						}
					}
				}
			},
		},

		///////////////////////////////////////////////////list//////////////////////////////////////////////////////
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

		///////////////////////////////////////////////////log///////////////////////////////////////////////////////
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
