package wordpress

import (
	"io/fs"
	"os"
	"path/filepath"
)


func CreateDirectories() string {
    if err := os.Mkdir("./dump", fs.FileMode(0755)); err != nil {
        panic(err)
    }

    if err := os.Mkdir("./wp", fs.FileMode(0755)); err != nil {
        panic(err)
    }

    result, _ := filepath.Abs("./wp")

    return result
}
