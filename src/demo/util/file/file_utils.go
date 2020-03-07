package file

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

func LoadFileContent(fileName string) []string {
	file, err := os.Open(GetParentPath(3) + "/" + fileName)
	if err != nil {
		panic(fmt.Sprintf("failed opening file: %s", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var content []string
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

func GetParentPath(level int) string {
	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parentPath := ""
	for i := 0; i < level; i++ {
		parentPath += "../"
	}
	dir, err := os.Open(path.Join(dirname, parentPath))
	if err != nil {
		panic(err)
	}
	return dir.Name()
}