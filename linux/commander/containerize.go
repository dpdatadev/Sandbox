package main

//UPGRADE GO TO 1.25
//If Docker is installed - the user can opt to run commands inside a container and redirect the output
//TODO
//https://github.com/ahmetb/go-dexec
/*
import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func dockerTest() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

*/
