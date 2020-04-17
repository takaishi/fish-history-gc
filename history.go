package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type Entry struct {
	Cmd  string `yaml:"cmd"`
	When int    `yaml:"when"`
}

type Entries []*Entry

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

func openFishHistory(path string) (_ *os.File, err error) {
	if path == "" {
		path, err = defaultHistoryPath()
		if err != nil {
			return nil, err
		}
	}

	return os.Open(path)
}

func defaultHistoryPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(user.HomeDir, ".local", "share", "fish", "fish_history"), nil
}
