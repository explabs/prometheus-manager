package routers

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"net/http"
)

func StopRoute(w http.ResponseWriter, r *http.Request){
	switch r.URL.Path {
	case "/stop/prometheus":
		StopContainer("prometheus")
		fmt.Println(w, "prometheus stopped")
	}
}

func StopContainer(containerName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStop(ctx, containerName, nil); err != nil {
	}
}

