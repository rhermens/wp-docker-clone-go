package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/rhermens/wp-docker-clone/docker"
	"github.com/rhermens/wp-docker-clone/pma"
	"github.com/rhermens/wp-docker-clone/remote"
	"github.com/rhermens/wp-docker-clone/wordpress"
)

func main() {
    helpFlag := flag.Bool("h", false, "Show help")
    remoteFlag := flag.Bool("remote", false, "Download wp-content")
    dumpFlag := flag.String("dumpfile", "", "MySQL Dump file")
    pmaFlag := flag.Bool("pma", true, "Use PhpMyAdmin")
    outFlag := flag.String("o", "", "Out directory")
    flag.Parse()

    if *helpFlag {
        flag.PrintDefaults()
        os.Exit(0)
    }

    if *outFlag == "" {
        flag.PrintDefaults()
        log.Fatal("-o is required")
    }

    if strings.Contains(*outFlag, "~") {
        usr, _ := user.Current()
        *outFlag = strings.ReplaceAll(*outFlag, "~", usr.HomeDir)
    }

    if _, err := os.Stat(*outFlag); os.IsNotExist(err) {
        os.Mkdir(*outFlag, fs.FileMode(0755))
    }

    if err := os.Chdir(*outFlag); err != nil {
        log.Fatalf("Cannot chdir to %s", *outFlag)
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

    log.Printf("Created project in %s. Happy hacking", *outFlag)
}

