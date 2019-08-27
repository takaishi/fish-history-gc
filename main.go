package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	Cmd  string `yaml:"cmd"`
	When int    `yaml:"when"`
}

type Entries []*Entry

func Run(path string, overWrite bool) error {

	path, err := getHistoryPath(path)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	oldEntries, err := readEntries(file)
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
			if entry.Cmd == newEntry.Cmd && newEntry.When < entry.When {
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

func readEntries(r io.Reader) (Entries, error) {
	scanner := bufio.NewScanner(r)
	entries := Entries{}
	for scanner.Scan() {
		entry := Entry{}
		line := scanner.Text()
		if strings.HasPrefix(line, "- cmd: ") {
			entry.Cmd = line[7:]
		}

		scanner.Scan()
		line = scanner.Text()
		if strings.HasPrefix(line, "  when: ") {
			entry.When, _ = strconv.Atoi(line[8:])
		}
		entries = append(entries, &entry)
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
