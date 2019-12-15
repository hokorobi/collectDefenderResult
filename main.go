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
		fmt.Fprintf(os.Stderr, "Usage: %s [Options][dir]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	var (
		fname   = flag.String("f", "hoge", "読み込むファイルの名前")
		lnum    = flag.Int("l", 4, "ファイルから何行目を取得するか")
		outfile = flag.String("o", "", "結果の出力先ファイルパス")
	)

	flag.Parse()

	// pp.Print(flag.Args())
	// pp.Print(flag.NArg())
	dir, err := getDir()
	if err != nil {
		return err
	}
	err = filepath.Walk(dir, execWalkFunc(*fname, *lnum, *outfile))
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
		return "", errors.New("Is not directory.")
	}
	return "", errors.New("An unexpected result.")
}

func execWalkFunc(fname string, lnum int, outfile string) filepath.WalkFunc {
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
			if i < lnum {
				continue
			}

			parentDir := filepath.Base(filepath.Dir(path))
			var w *bufio.Writer
			if outfile == "" {
				w = bufio.NewWriter(os.Stdout)
			} else {
				fpo, err := os.OpenFile(outfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				defer fpo.Close()

				w = bufio.NewWriter(fpo)
			}
			_, err := fmt.Fprintf(w, "%s\t%s\n", parentDir, s.Text())
			if err != nil {
				return err
			}
			err = w.Flush()
			if err != nil {
				return err
			}
			return nil

		}
		return nil
	}
}

// vim: ff=unix
