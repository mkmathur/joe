package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const joe string = `
 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄ 
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀█░█▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ 
      ▐░▌    ▐░▌       ▐░▌▐░▌          
      ▐░▌    ▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄ 
      ▐░▌    ▐░▌       ▐░▌▐░░░░░░░░░░░▌
      ▐░▌    ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀▀▀ 
      ▐░▌    ▐░▌       ▐░▌▐░▌          
 ▄▄▄▄▄█░▌    ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ 
▐░░░░░░░▌    ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀▀▀▀▀▀▀      ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀ 
`
const version string = "1.0.0"
const gitignoreUrl = "https://github.com/github/gitignore/archive/master.zip"
const dataDir string = ".joe-data"

var dataPath = path.Join(os.Getenv("HOME"), dataDir)

func findGitignores() (a map[string]string, err error) {
	_, err = ioutil.ReadDir(dataPath)
	if err != nil {
		return nil, err
	}

	filelist := make(map[string]string)
	filepath.Walk(dataPath, func(filepath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".gitignore") {
			name := strings.ToLower(strings.Replace(info.Name(), ".gitignore", "", 1))
			filelist[name] = filepath
		}
		return nil
	})
	return filelist, nil
}

func availableFiles() (a []string, err error) {
	gitignores, err := findGitignores()
	if err != nil {
		return nil, err
	}

	availableGitignores := []string{}
	for key, _ := range gitignores {
		availableGitignores = append(availableGitignores, key)
	}

	return availableGitignores, nil
}

func generate(args string) {
	names := strings.Split(args, ",")

	gitignores, err := findGitignores()
	if err != nil {
		log.Fatal(err)
	}

	notFound := []string{}
	output := ""
	for _, name := range names {
		if filepath, ok := gitignores[strings.ToLower(name)]; ok {
			bytes, err := ioutil.ReadFile(filepath)
			if err == nil {
				output += "#### " + name + " ####\n"
				output += string(bytes)
				continue
			}
		} else {
			notFound = append(notFound, name)
		}
	}

	if len(notFound) > 0 {
		fmt.Printf("Unsupported files: %s\n", strings.Join(notFound, ", "))
		fmt.Println("Run `joe ls` to see list of available gitignores.")
		output = ""
	}
	if len(output) > 0 {
		output = "#### joe made this: http://goel.io/joe\n" + output
	}
	fmt.Println(output)
}

func main() {
	app := cli.NewApp()
	app.Name = joe
	app.Usage = "generate .gitignore files from the command line"
	app.UsageText = "joe command [arguments...]"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"list"},
			Usage:   "list all available files",
			Action: func(c *cli.Context) {
				availableGitignores, err := availableFiles()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%d supported .gitignore files:\n", len(availableGitignores))
				sort.Strings(availableGitignores)
				fmt.Printf("%s\n", strings.Join(availableGitignores, ", "))
			},
		},
		{
			Name:    "u",
			Aliases: []string{"update"},
			Usage:   "update all available gitignore files",
			Action: func(c *cli.Context) {
				fmt.Println("Updating gitignore files..")
				err := RemoveContents(dataPath)
				if err != nil {
					log.Fatal(err)
				}
				err = DownloadFiles(gitignoreUrl, dataPath)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "g",
			Aliases: []string{"generate"},
			Usage:   "generate gitignore files",
			Action: func(c *cli.Context) {
				if c.NArg() != 1 {
					cli.ShowAppHelp(c)
				} else {
					generate(c.Args()[0])
				}
			},
		},
	}
	app.Run(os.Args)
}
