package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
  "bytes"
  "io/ioutil"
  // "reflect"
)

// yamlから読み込むための構造体を定義

type DotzConfigToml struct {
	Tracked tracked
}
type tracked struct{
  Files [][]interface{}
}

func ReadDotzConf(fileName string) (config DotzConfigToml){
  var configs DotzConfigToml
  _, err := toml.DecodeFile(fileName, &configs)
  if err != nil{
    fmt.Println(err)
  }

  return configs
}

func WriteDotzConf(config DotzConfigToml, filename string){
   if !FileExists(filename){
    return
  }

  var buff bytes.Buffer
  if err := toml.NewEncoder(&buff).Encode(config); err != nil{
    fmt.Println(err)
  }
  writeBuff := []byte(buff.String())
  err := ioutil.WriteFile(filename, writeBuff, 0644)
  if err != nil {
    fmt.Println(err)
  }
}

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
		HOME      string
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
			Name:        "HOME",
			Destination: &HOME,
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "push, p",
					Usage: "enable `git push`",
				},
			},
			Action: func(c *cli.Context) error {

				out, _ := exec.Command("git", "-C", DOTZ_ROOT, "status", "-z").Output()

				if len(out) != 0 {
					// stageに変更があるとき
					exec.Command("git", "-C", DOTZ_ROOT, "add", "-A").Run()

					commitMessage := "[dotz][backup] " + time.Now().Format("2006-01-02 15:04'05")
					fmt.Println(commitMessage)
					exec.Command("git", "-C", DOTZ_ROOT, "commit", "-m", commitMessage).Run()

				} else {
					fmt.Println("No change")
				}

				if c.Bool("push") {
					exec.Command("git", "-C", DOTZ_ROOT, "push").Run()
					fmt.Println("pushed")
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

		config := ReadDotzConf("./dotzconfig.toml")
		// fmt.Println(config)

    WriteDotzConf(config, DOTZ_ROOT+"dotz.toml")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
