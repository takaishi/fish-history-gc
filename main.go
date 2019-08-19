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

type Entries []Entry

func Run(path string, overWrite bool) error {

	path, err := getHistoryPath(path)
	if err != nil {
		return err
	}

	oldEntries, err := readEntries(path)
	if err != nil {
		return err
	}
	newEntries := removeDupEntries(oldEntries)

	sort.Slice(newEntries, func(i, j int) bool {
		return newEntries[i].When < newEntries[j].When
	})

	if overWrite {
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

func removeDupEntries(entries Entries) Entries {
	newEntries := Entries{}

	for _, entry := range entries {
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

	return newEntries
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

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] [/fish/history/path]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var overWrite bool

	flag.Usage = usage
	flag.BoolVar(&overWrite, "overwrite", false, "Overwrite entries")
	flag.Parse()

	path := flag.Arg(0)
	err := Run(path, overWrite)
	if err != nil {
		log.Fatal(err)
	}
}
