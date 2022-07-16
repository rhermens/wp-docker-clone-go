package main

import (
	"flag"
	"io/fs"
	"os"
	"os/user"
	"strings"

	"github.com/rhermens/wp-docker-clone/docker"
	"github.com/rhermens/wp-docker-clone/wordpress"
	"github.com/rhermens/wp-docker-clone/remote"
)

func main() {
    ftpFlag := flag.Bool("ftp", false, "Download wp-content")
    outFlag := flag.String("o", ".", "Out directory")

    flag.Parse()

    if strings.Contains(*outFlag, "~") {
        usr, _ := user.Current()
        *outFlag = strings.ReplaceAll(*outFlag, "~", usr.HomeDir)
    }

    if *ftpFlag {
        remote.DownloadThemesAndPlugins(*outFlag + "/wp")
    }

    if _, err := os.Stat(*outFlag); os.IsNotExist(err) {
        os.Mkdir(*outFlag, fs.FileMode(0755))
    }
    if err := os.Chdir(*outFlag); err != nil {
        panic(err)
    }

    composeFile := docker.NewDockerCompose()
    wordpress.AddToCompose(&composeFile)
    
    wpDir := wordpress.CreateDirectories()
    composeFile.Store()

    if *ftpFlag {
        remote.DownloadThemesAndPlugins(wpDir)
    }
}

