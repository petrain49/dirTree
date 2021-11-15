package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strconv"
)

type File struct {
	Info fs.DirEntry
	PathToFile string
}

type FileList []File

func (fl FileList) Len() int           { return len(fl) }
func (fl FileList) Swap(i, j int)      { fl[i], fl[j] = fl[j], fl[i] }
func (fl FileList) Less(i, j int) bool { return fl[i].Info.Name() < fl[j].Info.Name() }

func getContent(pathToFile string) (FileList, error) {
	if !fs.ValidPath(pathToFile) {
		return nil, errors.New("invalid path")
	}

	dir, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	list, err := dir.ReadDir(0)
	if err != nil {
		return nil, err
	}

	res := make(FileList, 0, len(list))

	for _, entry := range list {
		res = append(res, File{entry, path.Join(pathToFile, entry.Name())})
	}

	sort.Sort(res)

	return res, nil
}

func printFiles(files FileList, depth uint) error {
	var err error
	var name string
	var pathToFile string
	var fileInfo fs.FileInfo
	var size string
	var innerFiles FileList

	for x, file := range files {
		fileInfo, err = file.Info.Info()
		if err != nil {
			return err
		}

		name = fileInfo.Name()
		pathToFile = file.PathToFile

		size = strconv.FormatInt(fileInfo.Size(), 10) + "b"
		if size == "0b" {
			size = "empty"
		}

		printTab(depth, x == len(files)-1)

		if !fileInfo.IsDir() {
			fmt.Printf("%s (%s)\n", name, size)
		} else {
			fmt.Printf("%s\n", name)

			innerFiles, err = getContent(pathToFile)
			if err != nil {
				return err
			}

			err = printFiles(innerFiles, depth + 1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printTab(n uint, isLast bool) {
	a := "├───"
	b := "└───"
	//c := "│   "

	var x uint

	for x = 0; x < n; x++ {
		fmt.Print("    ")
	}
	if !isLast {
		fmt.Print(a)
	} else {
		fmt.Print(b)
	}
}