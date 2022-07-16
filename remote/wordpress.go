package remote

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type config struct {
    targetDir string
    sourceDir string
    username string
    password string
    host string
}

func setup(targetDir string) config {
    config := config {
        targetDir: targetDir,
    }

    fmt.Print("Hostname: ")
    fmt.Scanln(&config.host)

    fmt.Print("Username: ")
    fmt.Scanln(&config.username)

    fmt.Print("Password: ")
    fmt.Scanln(&config.password)

    fmt.Print("Remote directory: ")
    fmt.Scanln(&config.sourceDir)

    return config
}

func DownloadThemesAndPlugins(targetDir string) {
    cnf := setup(targetDir)

    c, err := ftp.Dial(cnf.host, ftp.DialWithTimeout(time.Second * 10))
    if err != nil {
        log.Fatal("Cant connect")
    }

    if err := c.Login(cnf.username, cnf.password); err != nil {
        log.Fatal("Cant login")
    }

    if cnf.sourceDir != "" {
        if err := c.ChangeDir(cnf.sourceDir); err != nil {
            log.Fatal("Cannot chdir")
        }
    }

    wpRoot, _ := c.CurrentDir()
    walker := c.Walk(wpRoot + "/wp-content")

    for walker.Next() {
        if walker.Stat().Type != ftp.EntryTypeFile {
            continue
        }

        sourceResponse, err := c.Retr(walker.Path())
        if err != nil {
            log.Fatal("sourceResponse err")
        }

        targetDir := cnf.targetDir + strings.Replace(
            walker.Path()[0:strings.LastIndex(walker.Path(), walker.Stat().Name)],
            wpRoot, 
            "", 
            1,
        )
        if _, err := os.Stat(targetDir); os.IsNotExist(err) {
            if err := os.MkdirAll(targetDir, fs.FileMode(0755)); err != nil {
                log.Fatal("Cant create target dir")
            }
        }

        file, err := os.Create(targetDir + walker.Stat().Name)
        if err != nil {
            log.Fatal("Cant create target file" + walker.Path())
        }

        source, _ := ioutil.ReadAll(sourceResponse)
        if _, err := file.Write(source); err != nil {
            log.Fatal("Cant write source")
        }

        file.Close()
        sourceResponse.Close()
    }

    c.Quit()
}
