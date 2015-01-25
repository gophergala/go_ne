package shared

import "fmt"

func DirectoryExists(directory string) string {
	return fmt.Sprintf("-d %v", directory)
}
