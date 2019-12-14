package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := _main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func _main() error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f name][-l num][dir]\n", os.Args[0])
		flag.PrintDefaults()
	}
	var (
		fname string
		lnum  int
	)
	flag.StringVar(&fname, "f", "hoge", "読み込むファイルの名前")
	flag.IntVar(&lnum, "l", 4, "ファイルから何行目を取得するか")

	flag.Parse()

	// pp.Print(flag.Args())
	// pp.Print(flag.NArg())
	dir, err := getDir()
	if err != nil {
		return err
	}
	err = filepath.Walk(dir, execWalkFunc(fname, lnum))
	if err != nil {
		return err
	}
	return nil
}

func getDir() (string, error) {
	var dir string
	if flag.NArg() < 1 {
		dir, _ = os.Getwd()
		return dir, nil
	}
	// pp.Print(flag.Args())
	dir = flag.Args()[0]
	f, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s is not exist.", dir)
	}
	if err != nil {
		return "", err
	}
	if !f.IsDir() {
		return "", errors.New("Is not exists.")
	}
	return dir, nil
}

func execWalkFunc(fname string, lnum int) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) != fname {
			return nil
		}
		fp, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fp.Close()

		s := bufio.NewScanner(fp)
		for i := 1; s.Scan(); i++ {
			if i == lnum {
				fmt.Printf("%s\t%s\n", filepath.Base(filepath.Dir(path)), s.Text())
				return nil
			}
		}
		return nil
	}
}

// vim: ff=unix
