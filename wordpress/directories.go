package wordpress

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)


func CreateDirectories() string {
    if err := os.Mkdir("./dump", fs.FileMode(0755)); err != nil {
        log.Fatal("Cannot mkdir dump")
    }

    if err := os.Mkdir("./wp", fs.FileMode(0755)); err != nil {
        log.Fatal("Cannot mkdir wp")
    }

    result, _ := filepath.Abs("./wp")

    return result
}
