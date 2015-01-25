package shared

import "fmt"

// DirectoryExists returns the condition to determine if the
// given directory exists in bash, e.g. -d directory
func DirectoryExists(directory string) string {
	return fmt.Sprintf("-d %v", directory)
}
