package readers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open file, err: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read data from file, err: %v", err)
	}

	return string(data), nil
}
func ReadOneLine() string {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		log.Fatal("failed to read line")
	}

	return line
}
