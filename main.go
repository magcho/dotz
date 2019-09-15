package main

import (
	"bytes"
	"errors"
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
	DOTZ_ROOT            string
	HOME                 string
	DOTZ_CONFIG_FILENAME string = "dotzconfig.toml"
	CONFIG               DotzConfigToml
)
var (
	ErrAlreadyFile              = errors.New("this file is already tracked")
	ErrInputFileTypeIsNotFonund = errors.New("input file type is not found")
	ErrFilesNotFound            = errors.New("files not found")
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
		return false
	} else {
		return true
	}
}

func FileExists(filename string) bool {
	if f, err := os.Stat(filename); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}
func createLink(filePath string, DOTZ_ROOT string, silentFlag bool) (dotzPath string, err error) {
	dotzSubDir, fileName := ParseFilePath(filePath)

	if FileExists(DOTZ_ROOT+dotzSubDir+fileName) || FolderExists(DOTZ_ROOT+dotzSubDir+fileName) {
		fmt.Println("> Tracked ", filePath)
		return "", ErrAlreadyFile
	}

	if FileExists(filePath) {
		// ファイルのとき
		if !FolderExists(DOTZ_ROOT + dotzSubDir) {
			// トラッキング対象がDOTZ_ROOTのサブフォルダにあるときに、サブフォルダを作成する
			exec.Command("mkdir", "-p", DOTZ_ROOT+dotzSubDir).Run()
		}
		exec.Command("mv", filePath, DOTZ_ROOT+dotzSubDir+fileName).Run()

		command := exec.Command("ln", "-sv", DOTZ_ROOT+dotzSubDir+fileName, HOME+dotzSubDir+fileName)

		if silentFlag {
			command.Run()
		} else {
			out, _ := command.Output()
			fmt.Printf("%s\n", out)
		}
		return DOTZ_ROOT + dotzSubDir + fileName, nil

	} else if FolderExists(filePath) {
		leefFolderName := fileName
		exec.Command("mv", filePath, DOTZ_ROOT+dotzSubDir+leefFolderName).Run()

		command := exec.Command("ln", "-sv", DOTZ_ROOT+dotzSubDir+leefFolderName, HOME+dotzSubDir+leefFolderName)

		if silentFlag {
			command.Run()
		} else {
			out, _ := command.Output()
			fmt.Printf("%s\n", out)
		}
		return DOTZ_ROOT + dotzSubDir + leefFolderName, nil

	} else {
		return "", ErrInputFileTypeIsNotFonund
	}
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

func convertFullPath(origin string) (string, error) {
	i, _ := exec.Command("pwd").Output()
	pwd := string(i[:len(i)-1]) // 末端に\nがあるので消してstringにキャスト
	dirArr := strings.Split(pwd, "/")

	if strings.HasPrefix(origin, HOME) {
		// フルパスが入力された時
		return origin, nil

	} else if strings.HasPrefix(origin, "../") {
		// 相対パスで上のディレクトリが指定されたとき
		n := strings.Count(origin, "../")
		basePath := strings.Join(dirArr[:len(dirArr)-n], "/")
		relativePath := strings.Replace(origin, "../", "", -1)
		return basePath + "/" + relativePath, nil

	} else if strings.HasPrefix(origin, "./") {
		// 相対パスでカレントディレクトリ以下が指定された時
		return strings.Replace(origin, "./", pwd+"/", 1), nil

	} else if FileExists(pwd+`/`+origin) || FolderExists(pwd+`/`+origin) {
		// 直接ファイル名が入力された時
		return pwd + "/" + origin, nil

	}
	return "", ErrInputFileTypeIsNotFonund
}

func main() {

	app := cli.NewApp()
	app.Name = "dotz"
	app.Version = "0.1.1"
	app.Usage = "Manage dotfiles and more"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "magcho",
			Email: "mail@magcho.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			// dotz project root
			Name:  "dotzRoot",
			Usage: "dotz project root path",
			// Value:       "~/.dotz/", // default
			Destination: &DOTZ_ROOT,
			EnvVar:      "DOTZ_ROOT",
		},
		cli.StringFlag{
			// 環境変数からの引き継ぎ
			Name:        "HOME",
			Usage:       "sysmtem user home directory. When especially needed only, Do this option",
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

		if FileExists(DOTZ_ROOT + DOTZ_CONFIG_FILENAME) {
			CONFIG = ReadDotzConf(DOTZ_ROOT + DOTZ_CONFIG_FILENAME)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Init dotz project",
			Action: func(c *cli.Context) error {

				if !FolderExists(DOTZ_ROOT) {
					exec.Command("mkdir", DOTZ_ROOT).Run()
				}
				if !FolderExists(DOTZ_ROOT + "/.git") {
					exec.Command("git", "-C", DOTZ_ROOT, "init").Run()
				}

				return nil
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "Restore dotz project and create symbric link",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "silent, s",
					Usage: "Hide ln log",
				},
			},
			Action: func(c *cli.Context) error {

				for _, item := range CONFIG.Tracked.Files {
					sLinkPath := replaceTilde2HomePath(item[0])
					entityPath := replaceSlash2DotzPath(item[1])

					if FileExists(entityPath) {
						sLinkSubDir, _ := ParseFilePath(sLinkPath)
						sLinkSubDir = HOME + sLinkSubDir
						if !FolderExists(sLinkSubDir) {
							exec.Command("mkdir", "-p", sLinkSubDir).Run()
						}

						command := exec.Command("ln", "-sv", entityPath, sLinkPath)
						if c.Bool("silent") {
							command.Run()
						} else {
							out, _ := command.Output()
							fmt.Printf("%s\n", out)
						}
					} else if FolderExists(entityPath) {
						if sLinkPath[len(sLinkPath)-1:len(sLinkPath)] == "/" {
							sLinkPath = sLinkPath[:len(sLinkPath)-1]
						}
						command := exec.Command("ln", "-sv", entityPath, sLinkPath)
						if c.Bool("silent") {
							command.Run()
						} else {
							out, _ := command.Output()
							fmt.Printf("%s\n", out)
						}
					} else {
						return ErrInputFileTypeIsNotFonund
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
					Usage: "Enable `git push`",
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
					fmt.Println("> No change")
				}

				if c.Bool("push") {
					exec.Command("git", "-C", DOTZ_ROOT, "push").Run()
					fmt.Println("> git pushed")
				}

				return nil
			},
		},
		{
			Name:    "track",
			Aliases: []string{"t"},
			Usage:   "Track files for dotz project",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "silent, s",
					Usage: "don't print link info",
				},
				cli.BoolFlag{
					Name:  "folder, f",
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

					originFilePath, err := convertFullPath(originFilePath)
					if err != nil {
						return err
					}

					if FileExists(originFilePath) {
						// 指定された対象がファイルの時
						dotzFilePath, err := createLink(originFilePath, DOTZ_ROOT, c.Bool("silent"))

						if err != ErrAlreadyFile {
							CONFIG.Tracked.Files = append(CONFIG.Tracked.Files, []string{
								replaceHomePath2Tilde(originFilePath),
								replaceDotzPath2Slash(dotzFilePath),
							})
						}

					} else if FolderExists(originFilePath) {
						// 指定された対象がフォルダの時

						// 引数で渡した時に最後に/をつけるかつけないかのゆらぎを統一
						// abc/def/ -> abc/def
						if originFilePath[len(originFilePath)-1:len(originFilePath)] == "/" {
							originFilePath = originFilePath[:len(originFilePath)-1]
						}

						if c.Bool("folder") {
							// fodler オプションを渡している
							dotzFolderPath, err := createLink(originFilePath, DOTZ_ROOT, c.Bool("silent"))

							if err != ErrAlreadyFile {
								CONFIG.Tracked.Files = append(CONFIG.Tracked.Files, []string{
									replaceHomePath2Tilde(originFilePath + "/"),
									replaceDotzPath2Slash(dotzFolderPath),
								})
							}
						} else {
							fmt.Println("Require -f option for folder tracking")
						}
					} else {
						return ErrFilesNotFound
					}
				}

				WriteDotzConf(CONFIG, DOTZ_ROOT+DOTZ_CONFIG_FILENAME)

				return nil
			},
		},
	}

	// app.Action = func(c *cli.Context) error {
	// 	return nil
	// }

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
