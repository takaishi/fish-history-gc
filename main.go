package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
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

func Action(c *cli.Context) error {
	pathFlg := c.String("path")

	path, err := getHistoryPath(pathFlg)
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

	err = writeEntries(path, newEntries)
	if err != nil {
		return err
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

func writeEntries(path string, entries Entries) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, entry := range entries {
		file.WriteString(fmt.Sprintf("- cmd: %s\n", entry.Cmd))
		file.WriteString(fmt.Sprintf("  when: %d\n", entry.When))
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
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "fish_history path",
		},
	}
	app.Action = func(c *cli.Context) error {
		return Action(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
