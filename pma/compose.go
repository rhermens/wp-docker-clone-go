package pma

import (
	"github.com/rhermens/wp-docker-clone/docker"
)

func AddToCompose(compose *docker.DockerCompose) {
    compose.Services["pma"] = docker.Service{
        Image: "phpmyadmin",
        Restart: "unless-stopped",
        Ports: []string { "8888:80" },
    }
}
