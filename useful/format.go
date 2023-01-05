package useful

import "strings"

func GetFilename(filepath string) string {
	strs := strings.Split(filepath, "/")

	filename := strs[len(strs)-1]

	return strings.Join(strings.Split(filename, "-")[1:], "-")
}
