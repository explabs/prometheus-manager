package routers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func StartContainer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageName := "prom/prometheus"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	pwd := os.Getenv("HOST_PWD")
	sourceConfigPath := path.Join(pwd, os.Getenv("CONFIG"))
	destConfigPath := "/etc/prometheus/prometheus.yml"
	sourceDataPath := path.Join(pwd, os.Getenv("DATA"))
	destDataPath := "/data"

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourceConfigPath,
				Target: destConfigPath,
			},
			{
				Type:   mount.TypeBind,
				Source: sourceDataPath,
				Target: destDataPath,
			},
		},
		AutoRemove: true,
	}, nil, nil, "prometheus")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "promtheus id: %s\n", resp.ID)
}
