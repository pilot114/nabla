package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/cheggaaa/pb/v3"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "add", Description: "add a task to the list"},
		{Text: "template", Description: "options for task template"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main2() {
	config, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Print(err)
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(config, &data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
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

}
