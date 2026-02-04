package pkg

import (
	"os/exec"
	"strings"
)

func Selector(fileName string) (path string, err error) {
	// Start Python script
	cmd := exec.Command("python3", fileName)
	output, err := cmd.Output()

	if err != nil {
		return "", err

	}
	return strings.TrimSpace(string(output)), nil
}
