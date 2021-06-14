// Package accio generates an initial set of templates to help you get started in
// developing your Go project.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/julieqiu/derrors"
)

var supportedScript = map[string]bool{
	"cmd":  true,
	"bash": true,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: accio [cmd|bash] [dir]")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	dir, err := filepath.Abs(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if script := flag.Arg(0); supportedScript[script] {
		err := createProjectDir(ctx, dir, flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
	log.Printf("Created %s script at %q", flag.Arg(0), dir)
}

func createProjectDir(ctx context.Context, dir string, script string) (err error) {
	defer derrors.WrapStack(&err, "createProjectDir(ctx, %q)", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
	files, err := ioutil.ReadDir(script)
	if err != nil {
		return err
	}
	for _, f := range files {
		inDir, err := filepath.Abs(script)
		if err != nil {
			return err
		}

		in, err := os.Open(fmt.Sprintf("%s/%s", inDir, f.Name()))
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(fmt.Sprintf("%s/%s", dir, f.Name()))
		if err != nil {
			return err
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()

		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		if err := out.Sync(); err != nil {
			return err
		}
	}

	tidyModule(dir)
	return err
}

func tidyModule(dir string) (err error) {
	defer derrors.WrapStack(&err, "tidyModule(%q)", dir)

	cmd := exec.Command("go", "mod", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
