package docker

import (
	"io/fs"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Service struct {
    Image string `yaml:"image"` 
    User string `yaml:"user"`
    Restart string `yaml:"restart"`
    Ports []string `yaml:"ports"`
    Environment map[string]string `yaml:"environment"`
    Volumes []string `yaml:"volumes"`
}

type Volume struct {

}

type DockerCompose struct {
    Version string `yaml:"version"`
    Services map[string]Service `yaml:"services"`
    Volumes map[string]Volume `yaml:"volumes"`
}

func (y DockerCompose) Store() {
    out, err := yaml.Marshal(y)
    if err != nil {
        panic(err)
    }

    err = os.WriteFile("docker-compose.yml", out, fs.FileMode(0644))
    if err != nil {
        log.Fatal("Cannot write docker-compose.yml")
    }
}

func NewDockerCompose() DockerCompose {
    return DockerCompose{
        Version: "3.1",
        Services: map[string]Service{},
        Volumes: map[string]Volume{},
    }
}

