package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
)

type Entry struct {
	Cmd  string `yaml:"cmd"`
	When int    `yaml:"when"`
}

type Entries []*Entry

func Run(path string, inPlace bool) error {

	path, err := getHistoryPath(path)
	if err != nil {
		return err
	}

	oldEntries, err := readEntries(path)
	if err != nil {
		return err
	}
	newEntries := Entries{}

	for _, entry := range oldEntries {

		exists := false
		for _, newEntry := range newEntries {
			if entry.Cmd == newEntry.Cmd {
				newEntry.When = entry.When
				exists = true
			}
		}
		if !exists {
			newEntries = append(newEntries, entry)
		}
	}

	sort.Slice(newEntries, func(i, j int) bool {
		return newEntries[i].When < newEntries[j].When
	})

	if inPlace {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		err = writeEntries(file, newEntries)
		if err != nil {
			return err
		}
	} else {
		err = writeEntries(os.Stdout, newEntries)
		if err != nil {
			return err
		}
	}

	return nil
}

func readEntries(path string) (Entries, error) {
	entries := Entries{}

	histBytes, err := ioutil.ReadFile(path)

	err = yaml.Unmarshal(histBytes, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func writeEntries(io io.Writer, entries Entries) error {
	for _, entry := range entries {
		io.Write([]byte(fmt.Sprintf("- cmd: %s\n", entry.Cmd)))
		io.Write([]byte(fmt.Sprintf("  when: %d\n", entry.When)))
	}

	return nil
}

func getHistoryPath(pathFlag string) (string, error) {
	if pathFlag == "" {
		return defaultHistoryPath()
	} else {
		return pathFlag, nil
	}
}

func defaultHistoryPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(user.HomeDir, ".local", "share", "fish", "fish_history"), nil
}

func main() {
	var inPlace bool
	flag.BoolVar(&inPlace, "in-place", false, "aaa")
	flag.BoolVar(&inPlace, "s", false, "aaa")
	flag.Parse()

	path := flag.Arg(0)
	err := Run(path, inPlace)
	if err != nil {
		log.Fatal(err)
	}
}
