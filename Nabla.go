package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

func readConfig(fileName string) map[string]interface{} {
	config, err := ioutil.ReadFile(fileName)
	if err != nil {
	fmt.Print(err)
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(config, &data)
	if err != nil {
	log.Fatalf("error: %v", err)
	}
	return data;
}

func readAndWatchConfig() *map[string]interface{} {

	config := readConfig("config.yaml");

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
					config = readConfig(event.FileInfo.Name());
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

	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()

	return &config;
}

func getKeys(config *map[string]interface{}) []string {
	keys := make([]string, len(*config))
	i := 0
	for k := range *config {
		keys[i] = k
		i++
	}
	return keys
}

/**
1. Чтение из yaml набора команд, файл вотчится
2. два режима - запуск с аргементами/флагами или интерактивный
3. промпт будет помогать вводить команды и опции
4. реестр хэндлеров
5. некоторые команды могут вызывать хэндлеры для генерации подсказок
	*/
func main() {

	config := readAndWatchConfig();
	fmt.Println(config)

	for {
		keys := getKeys(config)
		fmt.Println(keys)
		time.Sleep(time.Second)
	}
}
