package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"path"
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
		os.Exit(-1)
	}

	templatePath := os.Args[1]
	configPath := os.Args[2]

	os.Exit(main_(templatePath, configPath))
}

func main_(templatePath string, configPath string) int {

	funcMap := template.FuncMap{
		"tfStringListFormater": tfStringListFormater,
	}

	templateName := path.Base(templatePath)
	template, err := template.New(templateName).Funcs(funcMap).ParseFiles(templatePath)

	if err != nil {
		fmt.Printf("error %v\n", err)
		return -1
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return -1
	}

	var context interface{}
	err = yaml.Unmarshal(data, &context)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return -1
	}

	err = template.Execute(os.Stdout, context)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return -1
	}
	return 0
}

func usage() {
	fmt.Printf("usage  : terraformer [go template file] [config yaml file]\n")
	fmt.Printf("version: %v \n", version)
	fmt.Printf("commit : %v \n", commit)
	fmt.Printf("date   : %v \n", date)
}

func tfStringListFormater(a [] interface{}) string {
	var result string = "[]"

	if a == nil {return result}
	if len(a) == 0 {return result}
	if len(a) == 1 {return fmt.Sprintf("[\"%v\"]", a[0])}

	result = "["
	for idx, val := range a {
		if idx < (len(a) - 1) {
			result = fmt.Sprintf("%v\"%v\", ", result, val)
		} else {
			result = fmt.Sprintf("%v\"%v\"]",  result, val)
		}
	}

	return result
}

