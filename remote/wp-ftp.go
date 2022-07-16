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
    connection *ftp.ServerConn
}
func (cnf *config) Close() {
    if err := cnf.connection.Quit(); err != nil {
        log.Fatal(err)
    }
}

func Setup(targetDir string) config {
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

    c, err := ftp.Dial(config.host, ftp.DialWithTimeout(time.Second * 10))
    if err != nil {
        log.Fatal("Cant connect")
    }

    if err := c.Login(config.username, config.password); err != nil {
        log.Fatal("Cant login")
    }

    if config.sourceDir != "" {
        if err := c.ChangeDir(config.sourceDir); err != nil {
            log.Fatal("Cannot chdir")
        }
    }

    config.connection = c

    return config
}

func DownloadThemesAndPlugins(cnf *config) {
    wpRoot, _ := cnf.connection.CurrentDir()
    walker := cnf.connection.Walk(wpRoot + "/wp-content")

    for walker.Next() {
        if walker.Stat().Type != ftp.EntryTypeFile {
            continue
        }

        log.Printf("Downloading: %s\n", walker.Path())

        sourceResponse, err := cnf.connection.Retr(walker.Path())
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
}

func DownloadDumpFile(cnf *config, dumpPath string) {
    sourceResponse, err := cnf.connection.Retr(dumpPath)
    if err != nil {
        log.Fatal("sourceResponse err")
    }

    file, err := os.Create(cnf.targetDir + "/dump.sql")
    if err != nil {
        log.Fatal("Cant create dump file")
    }

    source, _ := ioutil.ReadAll(sourceResponse)
    if _, err := file.Write(source); err != nil {
        log.Fatal("Cant write source")
    }

    file.Close()
    sourceResponse.Close()
}
