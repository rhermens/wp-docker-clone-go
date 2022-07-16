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
	"github.com/rhermens/wp-docker-clone/pma"
)

func main() {
    remoteFlag := flag.Bool("remote", false, "Download wp-content")
    dumpFlag := flag.String("dumpfile", "", "MySQL Dump file")
    pmaFlag := flag.Bool("pma", true, "Use PhpMyAdmin")
    outFlag := flag.String("o", ".", "Out directory")
    flag.Parse()

    if strings.Contains(*outFlag, "~") {
        usr, _ := user.Current()
        *outFlag = strings.ReplaceAll(*outFlag, "~", usr.HomeDir)
    }

    if _, err := os.Stat(*outFlag); os.IsNotExist(err) {
        os.Mkdir(*outFlag, fs.FileMode(0755))
    }

    if err := os.Chdir(*outFlag); err != nil {
        panic(err)
    }

    composeFile := docker.NewDockerCompose()
    wordpress.AddToCompose(&composeFile)
    if *pmaFlag {
        pma.AddToCompose(&composeFile) 
    }
    composeFile.Store()

    wpDir := wordpress.CreateDirectories()
    if *remoteFlag {
        ftpCnf := remote.Setup(wpDir)
        remote.DownloadThemesAndPlugins(&ftpCnf)

        if *dumpFlag != "" {
            remote.DownloadDumpFile(&ftpCnf, *dumpFlag)
        }

        ftpCnf.Close()
    }
}

