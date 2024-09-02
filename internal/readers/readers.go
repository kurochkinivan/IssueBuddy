package readers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	cuserr "github.com/kurochkinivan/IssueBuddy/internal/custome_erros"
)

func ReadFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("%s %v", cuserr.OpenFileErr, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s %v", cuserr.ReadFileErr, err)
	}

	return string(data), nil
}
func ReadOneLine() (string, error) {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read line, err: %v", err)
	}

	return strings.TrimSpace(line), nil
}
