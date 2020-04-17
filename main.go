package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	var overWrite bool

	flag.Usage = usage
	flag.BoolVar(&overWrite, "overwrite", false, "Overwrite entries")
	flag.Parse()

	path := flag.Arg(0)
	if err := Run(path, overWrite); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] [/fish/history/path]\n", os.Args[0])
	flag.PrintDefaults()
}

func Run(path string, overWrite bool) error {
	lock := NewFileLock("/tmp/fish-history-gc.lock")
	lock.Lock()
	defer lock.Unlock()

	file, err := openFishHistory(path)
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
