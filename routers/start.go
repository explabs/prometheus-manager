package routers

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type Config struct {
	Services []Service `yaml:"services"`
}
type Service struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Ports   []string `yaml:"ports"`
	Volumes []string `yaml:"volumes"`
}

func (cfg *Config) LoadConfig(configName string) error {
	f, err := os.Open(configName)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SearchByName(name string) *Service {
	for _, service := range cfg.Services {
		if service.Name == name {
			return &service
		}
	}
	return nil
}

func StartRoute(w http.ResponseWriter, r *http.Request) {
	var c Config
	c.LoadConfig("manager.yml")
	switch r.URL.Path {
	case "/start/prometheus":

		StartContainer(c.SearchByName("prometheus"))
		fmt.Println(w, "prometheus started")
	case "/start/malwaretotal":
		StartContainer(c.SearchByName("malwaretotal"))
		fmt.Println(w, "malwaretotal started")

	}
}

func CheckEnvs(envs []string) error {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			return fmt.Errorf("%s is empty", env)
		}
	}
	return nil
}

func StartContainer(s *Service) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, s.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	pwd := os.Getenv("PWD")
	if pwd == "" {
		return fmt.Errorf("PWD is empty")
	}

	var mounts []mount.Mount
	for _, volume := range s.Volumes {
		array := strings.Split(volume, ":")
		source := path.Join(pwd, array[0])
		target := array[1]
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		})
	}
	var ports nat.PortMap
	for _, port := range s.Ports {
		array := strings.Split(port, ":")
		source := array[0]
		target := array[1]
		ports = nat.PortMap{nat.Port(source + "/" + "tcp"): []nat.PortBinding{{HostPort: target}}}
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: s.Image,
	}, &container.HostConfig{
		Mounts:       mounts,
		AutoRemove:   true,
		PortBindings: ports,
	}, nil, nil, s.Name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
