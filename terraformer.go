package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}

	templatePath := os.Args[1]
	configPath := os.Args[2]

	template, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
	}

	var context interface{}
	yaml.Unmarshal(data, &context)

	template.Execute(os.Stdout, context)
}

func usage() {
	fmt.Printf("usage  : terraformer [go template file] [config yaml file]\n")
	fmt.Printf("version: %v \n", version)
	fmt.Printf("commit : %v \n", commit)
	fmt.Printf("date   : %v \n", date)
}
