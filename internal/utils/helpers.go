package utils

import "fmt"

func FormatMemory(memory int) string {
	return fmt.Sprintf("%dMiB", memory)
}
