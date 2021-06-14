package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/julieqiu/derrors"
)

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
	switch flag.Arg(0) {
	case "bash":
		err = accioBash(ctx, dir)
	case "cmd":
		err = accioCmd(ctx, dir)
	default:
		flag.Usage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created %s script at %q", flag.Arg(0), dir)
}

func accioBash(ctx context.Context, dir string) (err error) {
	defer derrors.WrapStack(&err, "accioBash(ctx, %q)", dir)

	if err := createProjectDir(ctx, dir); err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/main.go", dir))
	if err != nil {
		return err
	}
	f.WriteString(`package main

import (
	"bytes"
    "fmt"
    "log"
    "os/exec"
    "strings"
)

func main() {
    cmd := exec.Command("tr", "a-z", "A-Z")
    cmd.Stdin = strings.NewReader("some input")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(out.String())
}`)
	f.Sync()
	tidyModule(dir)
	return nil
}

func accioCmd(ctx context.Context, dir string) (err error) {
	defer derrors.WrapStack(&err, "accioCmd(ctx, %q)", dir)

	if err := createProjectDir(ctx, dir); err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/main.go", dir))
	if err != nil {
		return err
	}
	defer func() {
		cerr := f.Close()
		if err != nil {
			err = cerr
		}
	}()

	f.WriteString(fmt.Sprintf(`package main

import (
    "flag"
    "fmt"
)

// var myFlag = flag.Bool("flagname", false, "TODO...")

func main() {
    flag.Usage = func() {
        fmt.Fprintln(flag.CommandLine.Output(), "usage: %s [TODO(fill this in)]")
        flag.PrintDefaults()
    }

    flag.Parse()

    // if flag.NArg() != 2 {
    // Uncomment to check number of args.
    // }

    // switch flag.Arg(0) {
    // Uncomment to switch on the first arg.
    // }
}`, filepath.Base(dir)))
	tidyModule(dir)
	return nil
}

func createProjectDir(ctx context.Context, dir string) (err error) {
	defer derrors.WrapStack(&err, "createProjectDir(ctx, %q)", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
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
