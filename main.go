package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	// parse args
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("./update-compose compose-file image-name new-version")
		os.Exit(-1)
	}
	file := args[0]
	image := args[1]
	newVersion := args[2]

	// parse file
	f, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var compose Compose
	err = yaml.Unmarshal(f, &compose)
	if err != nil {
		panic(err)
	}

	// search for desired image
	var originals []string
	for _, s := range compose.Services {
		img := s.Image
		if !strings.Contains(img, ":") {
			continue
		}
		toks := strings.Split(img, ":")
		name := toks[0]
		if name != image {
			continue
		}
		originals = append(originals, img)
	}
	if len(originals) == 0 {
		log.Printf("Image: %q not found!", image)
		os.Exit(-1)
	}

	// replacement(s)
	cf := string(f)
	update := fmt.Sprintf("%s:%s", image, newVersion)
	for _, o := range originals {
		cf = strings.Replace(cf, o, update, -1)
	}

	// write result
	err = ioutil.WriteFile(file, []byte(cf), 0644)
	if err != nil {
		panic(err)
	}
	for _, o := range originals {
		if o == update {
			continue
		}
		log.Printf("Updated: %q => %q\n", o, update)
	}
}

// Compose represents the compose file.
type Compose struct {
	Services map[string]Service
}

// Service represents the service section.
type Service struct {
	Image string `yaml:"image"`
}
