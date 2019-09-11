package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func FolderExists(filename string) bool {
	if f, err := os.Stat(filename); os.IsNotExist(err) || !f.IsDir() {
		// fmt.Println("ディレクトリは存在しません！")
		return false
	} else {
		// fmt.Println("存在します")
		return true
	}
}

func FileExists(filename string) bool {
	if f, err := os.Stat(filename); os.IsNotExist(err) || f.IsDir() {
		// fmt.Println("ファイルは存在しません！")
		return false
	} else {
		// fmt.Println("存在するファイルです")
		return true
	}
}

func ParseFilePath(path string) (string, string) {
	r := regexp.MustCompile(`\s*/\s*`)
	result := r.Split(path, -1)
	fileName := result[len(result)-1]
	dirPath := strings.Join(result[:len(result)-2], "/") + "/"
	return dirPath, fileName
}

func main() {
	var (
		DOTZ_ROOT string
		home      string
	)

	app := cli.NewApp()
	app.Name = "dotz"
	app.Version = "0.1.0"
	app.Usage = "macOS backup and restore dotfiles"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "magcho",
			Email: "mail@magcho.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			// dotz project root
			Name:  "dotzRoot, r",
			Usage: "project root path",
			// Value:       "~/.dotz/", // default
			Destination: &DOTZ_ROOT,
			EnvVar:      "DOTZ_ROOT",
		},
		cli.StringFlag{
			// 環境変数からの引き継ぎ
			Name:        "home",
			Destination: &home,
			EnvVar:      "HOME",
		},
	}
	app.Before = func(c *cli.Context) error {

		// DOTZ_ROOTが"/"で終わってないときは補完
		// /hoge/fuga -> /hoge/fuga/
		if DOTZ_ROOT[len(DOTZ_ROOT)-1:len(DOTZ_ROOT)] != "/" {
			DOTZ_ROOT = DOTZ_ROOT + "/"
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init dotz project",
			Action: func(c *cli.Context) error {

				// 引数がない時
				// if c.NArg() < 1 {
				// 	fmt.Println("require dotz project name")
				// 	fmt.Println("dotz [init|i]")
				// 	return nil
				// }

				if !FolderExists(DOTZ_ROOT) {
					exec.Command("mkdir", DOTZ_ROOT).Run()
					fmt.Println("> mkdir " + DOTZ_ROOT)
				}
				if !FolderExists(DOTZ_ROOT + "/.git") {
					exec.Command("git", "-C", DOTZ_ROOT, "init").Run()
					fmt.Println("> git init")
				}

				return nil
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "restore dotz project and make symbric link",
			Action: func(c *cli.Context) error {

				fmt.Println()
				return nil
			},
		},
		{
			Name:    "backup",
			Aliases: []string{"b"},
			Usage:   "backup tracked files",
			Action: func(c *cli.Context) error {

				if c.NArg() < 1 {
					fmt.Println("require ")
				}
				for i := 0; i <= c.NArg(); i++ {
					fmt.Println(c.Args().Get(i))
				}
				return nil
			},
		},
		{
			Name:    "track",
			Aliases: []string{"t"},
			Usage:   "file append into dotz project",
			Action: func(c *cli.Context) error {

				// 引数がないとき
				if c.NArg() < 1 {
					fmt.Println("require file in args")
					return nil
				}

				for i := 0; i < c.NArg(); i++ {
					filePath := c.Args().Get(i)

					_, fileName := ParseFilePath(filePath)
					exec.Command("mv", filePath, DOTZ_ROOT).Run()

					out, _ := exec.Command("ln", "-sv", DOTZ_ROOT+fileName, filePath).Output()
					fmt.Printf("%s\n", out)
				}

				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		// var out_ string

		out, _ := exec.Command("pwd").Output()
		fmt.Printf("%s\n", out)

		out, _ = exec.Command("cd", "../").Output()
		fmt.Printf("%s\n", out)

		out, _ = exec.Command("cd", "../", ";", "pwd").Output()
		fmt.Printf("%s\n", out)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
