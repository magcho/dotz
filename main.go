package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	DOTZ_ROOT string
	HOME      string
	CONFIG    DotzConfigToml
)

// yamlから読み込むための構造体を定義
type DotzConfigToml struct {
	Tracked tracked
}
type tracked struct {
	// Files [][]interface{}
	Files [][]string
}

func ReadDotzConf(fileName string) (config DotzConfigToml) {
	var configs DotzConfigToml
	_, err := toml.DecodeFile(fileName, &configs)
	if err != nil {
		fmt.Println(err)
	}

	return configs
}

func WriteDotzConf(config DotzConfigToml, filePath string) {

	var buff bytes.Buffer
	if err := toml.NewEncoder(&buff).Encode(config); err != nil {
		fmt.Println(err)
	}
	writeBuff := []byte(buff.String())
	err := ioutil.WriteFile(filePath, writeBuff, 0644)
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
func createLink(filePath string, DOTZ_ROOT string, silentFlag bool) (dotzPath string, err error) {
	dotzSubDir, fileName := ParseFilePath(filePath)

  
	if FileExists(DOTZ_ROOT + dotzSubDir + fileName) || FolderExists(DOTZ_ROOT + dotzSubDir + fileName){
		fmt.Println(">", filePath, "this is already tracked")
		return
	}

  if FileExists(filePath){
    // ファイルのとき
    if !FolderExists(DOTZ_ROOT + dotzSubDir) {
      // トラッキング対象がDOTZ_ROOTのサブフォルダにあるときに、サブフォルダを作成する
      exec.Command("mkdir", "-p", DOTZ_ROOT + dotzSubDir).Run()
    }
    exec.Command("mv", filePath, DOTZ_ROOT + dotzSubDir + fileName).Run()

    command := exec.Command("ln", "-sv", DOTZ_ROOT + dotzSubDir + fileName, HOME + dotzSubDir + fileName)
    if !silentFlag {
      out, _ := command.Output()
      fmt.Printf("%s\n", out)
    } else {
      command.Run()
    }
    return DOTZ_ROOT + dotzSubDir + fileName, nil
  }else if FolderExists(filePath){
    leefFolderName := fileName
    exec.Command("mv", filePath, DOTZ_ROOT + dotzSubDir + leefFolderName).Run()
    
    
    command := exec.Command("ln", "-sv", DOTZ_ROOT + dotzSubDir + leefFolderName, HOME + dotzSubDir + leefFolderName)
    
    if silentFlag {
      command.Run()
    }else{
      out, _ := command.Output()
      fmt.Printf("%s\n", out)
    }
    return 
  }
  return
}

func ParseFilePath(path string) (string, string) {
	path = strings.Replace(path, HOME, "", 1)

	r := regexp.MustCompile(`\s*/\s*`)

	result := r.Split(path, -1)
	fileName := result[len(result)-1]

	subDirPath := strings.Replace(path, fileName, "", 1)
	return subDirPath, fileName
}

func replaceHomePath2Tilde(origin string) string {
	return strings.Replace(origin, HOME, "~/", 1)
}
func replaceTilde2HomePath(origin string) string {
	return strings.Replace(origin, "~/", HOME, 1)
}
func replaceDotzPath2Slash(origin string) string {
	return strings.Replace(origin, DOTZ_ROOT, "//", 1)
}
func replaceSlash2DotzPath(origin string) string {
	return strings.Replace(origin, "//", DOTZ_ROOT, 1)
}

func main() {

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
		if HOME[len(HOME)-1:len(HOME)] != "/" {
			HOME = HOME + "/"
		}

		if FileExists(DOTZ_ROOT + "dotzconfig.toml") {
			CONFIG = ReadDotzConf(DOTZ_ROOT + "dotzconfig.toml")
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
      Flags: []cli.Flag{
        cli.BoolFlag{
          Name: "silent, s",
          Usage: "don't print ln info",
        },
      },
			Action: func(c *cli.Context) error {

        for _, item := range CONFIG.Tracked.Files{
          fmt.Println(item[0], item[1])
          fmt.Println(replaceTilde2HomePath(item[0]), replaceSlash2DotzPath(item[1]))

          command := exec.Command("ln","-sv")
          if c.Bool("silent"){
            command.Run()
          }else{
            out, _ := command.Output()
            fmt.Println(out)
          }
        }
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "silent, s",
					Usage: "don't print link info",
				},
        cli.BoolFlag{
          Name: "folder, f",
          Usage: "manage in dotz of folder",
        },
			},
			Action: func(c *cli.Context) error {

				// 引数がないとき
				if c.NArg() < 1 {
					fmt.Println("require file in args")
					return nil
				}

       
				for i := 0; i < c.NArg(); i++ {
					originFilePath := c.Args().Get(i)
          
          if FileExists(originFilePath){
            // 指定された対象がファイルの時
            dotzFilePath, _ := createLink(originFilePath, DOTZ_ROOT, !c.Bool("silent"))

            CONFIG.Tracked.Files = append(CONFIG.Tracked.Files, []string{
              replaceHomePath2Tilde(originFilePath),
              replaceDotzPath2Slash(dotzFilePath),
            })
          }else if FolderExists(originFilePath){
            // 指定された対象がフォルダの時
            
            if c.Bool("folder"){
              // fodler オプションを渡している
              createLink(originFilePath, DOTZ_ROOT, c.Bool("silent"))
            }
          }else{
            fmt.Println("err")
          }
				}

				WriteDotzConf(CONFIG, DOTZ_ROOT+"dotzconfig.toml")
				return nil
			},
		},
	}




  app.Action = func(c *cli.Context) error {
    fmt.Println(FolderExists("/Users/magcho/sample/sample.txt/"))
    return nil
  }


  
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
