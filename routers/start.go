package routers

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"net/http"
	"os"
	"path"
)

var DestConfigPath = "/etc/prometheus/prometheus.yml"
var DestDataPath = "/data"

func StartRoute(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/start/prometheus":
		StartContainer("prometheus", "prom/prometheus", "9091", os.Getenv("HOST_PWD"), os.Getenv("CONFIG"), os.Getenv("DATA"))
		fmt.Println(w, "prometheus started")
	// add case for another service

	}
}

func StartContainer(containerName string, imageName string, port string, pwd string, config string, data string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	sourceConfigPath := path.Join(pwd, config)
	sourceDataPath := path.Join(pwd, data)

	resp, err := cli.ContainerCreate(ctx, &container.Config{

		Image: imageName,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourceConfigPath,
				Target: DestConfigPath,
			},
			{
				Type:   mount.TypeBind,
				Source: sourceDataPath,
				Target: DestDataPath,
			},
		},
		AutoRemove: true,
		PortBindings: nat.PortMap{
			nat.Port(port + "/" + "tcp"): []nat.PortBinding{{HostPort: port}},
		},
	}, nil, nil, containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}
