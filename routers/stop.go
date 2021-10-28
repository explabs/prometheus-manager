package routers

import (
	"context"
	"log"
	"net/http"

	"github.com/docker/docker/client"
)

func StopContainert(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerName := "prometheus"
	if err := cli.ContainerStop(ctx, containerName, nil); err != nil {
		log.Printf("Unable to stop container %s: %s", containerName, err)
	}
}
