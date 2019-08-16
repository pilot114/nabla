package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/cheggaaa/pb/v3"
	"github.com/urfave/cli"
	"github.com/radovskyb/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

type Command struct {
	Open string `yaml:"open"`
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "add", Description: "add a task to the list"},
		{Text: "template", Description: "options for task template"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

/**
1. Чтение из yaml набора команд, файл вотчится
2. два режима - запуск с аргементами/флагами или интерактивный
3. промпт будет помогать вводить команды и опции
4. реестр хэндлеров
5. некоторые команды могут вызывать хэндлеры для генерации подсказок
	*/
func main() {
	config, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Print(err)
	}

	com := Command{}
	err = yaml.Unmarshal(config, &com)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", com)
	fmt.Println(string(config))

	var cc string
	var b string

	app := cli.NewApp()
	// app -som "Some message"
	app.UseShortOptionHandling = true

	app.Name = "Nabla"
	app.Version = "0.2"
	app.Usage = "Выполнятор команд"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &cc,
			// взять из переменной окружения
			//EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
		},
		cli.StringFlag{
			Name:        "b, bbb",
			Value:       "fuzz",
			Usage:       "second flag",
			Destination: &b,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "add",
			Aliases:  []string{"a"},
			Usage:    "add a task to the list",
			Category: "Template actions",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:     "template",
			Aliases:  []string{"t"},
			Usage:    "options for task templates",
			Category: "Template actions",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q\n", c.Args().Get(0))
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	count := 1000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)

	w := watcher.New()
	w.SetMaxEvents(10)
	w.FilterOps(watcher.Write)

	r := regexp.MustCompile("^config.yaml$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {
					fmt.Println(event.FileInfo.Name())
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add("."); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
