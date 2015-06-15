package main

import (
	//"bytes"
	//"fmt"
	"github.com/dalanlan/dockerclient/docker"
	"log"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	//	endpoint := "http://127.0.0.1:4243"
	client := docker.NewClient(endpoint)

	/*log.Println("list images on your host")

	imgs, err := client.ListImage()
	if err != nil {
		log.Fatal(err)
	}
	for _, img := range imgs {
		fmt.Println("ID:", img.Id)
		fmt.Println("RepoTags:", img.RepoTags)
		fmt.Println("Create:", img.Created)
		fmt.Println("Size:", img.Size)
		fmt.Println("VirtualSize:", img.VirtualSize)
		fmt.Println("ParentId:", img.ParentId)
	}*/

	/*--------------------------------------------*/
	/*log.Println("list running containers on your host")
	containers, err := client.ListContainers()
	if err != nil {
		log.Fatal(err)
	}
	for _, cont := range containers {
		fmt.Println("ID:", cont.Id)
		fmt.Println("Ports:", cont.Ports)
	}*/

	/*--------------------------------------------*/
	/*	log.Println("create a container")

		config := &docker.Config{
			AttachStdin:  false,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    false,
			StdinOnce:    false,
			Env:          nil,
			Cmd:          []string{"date"},
			Image:        "ubuntu",
			Labels: map[string]string{
				"com.example.vendor":  "Acme",
				"com.example.license": "GPL",
				"com.example.version": "1.0",
			},
			/*Volumes: map[string]struct{}{
				"/tmp": {},
			},
			NetworkDisabled: false,
			MacAddress:      "12:34:56:78:9a:bc",
			ExposedPorts:map[&docker.Port]struct{}{
			}
		}*/

	/*hostConfig := &docker.HostConfig{}*/

	/*opts := docker.CreateContainerOption{
		Name:   "TestCreateContainer",
		Config: config,
	}
	container, err := client.CreateContainers(opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID", container.ID)*/

	/*--------------------------------------------*/
	/*log.Println("get container logs")
		var bufout bytes.Buffer
		var buferr bytes.Buffer
		logsopt := docker.GetContainerLogOption{
			Follow:     true,
			Stdout:     true,
			Stderr:     true,
			Timestamps: true,
			Tail:       "10",
			Container:  "3c5f38f3af27",
			OutStream:  &bufout,
			ErrStream:  &buferr,
		}

		if err = client.GetContainerLogs(logsopt); err != nil {
			log.Fatal(err)
		}

	fmt.Println("Stdout:", bufout.String())
		fmt.Println("Stderr:", buferr.String())*/

	/*--------------------------------------------*/
	log.Println("stop a container")
	stopOpt := docker.StopContainerOption{
		Time:      1000,
		Container: "3c5f38f3af27",
	}
	if err := client.StopContainer(stopOpt); err != nil {
		log.Fatal(err)
	}

}
