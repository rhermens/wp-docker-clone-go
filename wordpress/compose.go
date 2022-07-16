package wordpress

import (
	"github.com/rhermens/wp-docker-clone/docker"
)

func AddToCompose(compose *docker.DockerCompose) {
    addMysql(compose)

    compose.Services["wordpress"] = docker.Service{
        Image: "wordpress",
        User: "1000:1000",
        Restart: "unless-stopped",
        Ports: []string { "8080:80" },
        Environment: map[string]string{
            "WORDPRESS_DB_HOST": "db",
            "WORDPRESS_DB_USER": "wp_dev",
            "WORDPRESS_DB_PASSWORD": "wp_dev",
            "WORDPRESS_DB_NAME": "wp_dev",
            "WORDPRESS_CONFIG_EXTRA": "\ndefine('WP_HOME', 'http://localhost:8080');\ndefine('WP_SITEURL', 'http://localhost:8080');",
        },
        Volumes: []string{ "./wp:/var/www/html" },
    }
}

func addMysql(compose *docker.DockerCompose) {
    compose.Services["db"] = docker.Service{
        Image: "mysql:5.7",
        Restart: "unless-stopped",
        Environment: map[string]string{
            "MYSQL_DATABASE": "wp_dev",
            "MYSQL_USER": "wp_dev",
            "MYSQL_PASSWORD": "wp_dev",
            "MYSQL_RANDOM_ROOT_PASSWORD": "1",
        },
        Volumes: []string {
            "./dump:/docker-entrypoint-initdb.d",
            "db:/var/lib/mysql",
        },
    }

    compose.Volumes["db"] = docker.Volume{}
}
