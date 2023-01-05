package useful

import (
	"os"
)

func EnsureDirectories() {
	os.Mkdir("./uploads", os.ModeDir)
}
