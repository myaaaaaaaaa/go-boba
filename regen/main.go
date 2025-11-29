package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	readme := string(must(os.ReadFile("README.md")))
	readme, _, _ = strings.Cut(readme, "## Demo")
	readme += "## Demos\n\n"

	for _, gif := range must(filepath.Glob("*.gif")) {
		readme += fmt.Sprintf("#### %s\n![%s](%s)\n\n", gif, gif, gif)

	}

	must(0, os.WriteFile("README.md", []byte(readme), 0666))
}
