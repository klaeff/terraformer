package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"gopkg.in/alecthomas/kingpin.v2"
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
	app := kingpin.New("terraformer", "A go program that generates terraform files using go templates")
	app.Version(printVersion()).Author("Stephan Klevenz")

	commandGenerate := app.Command("generate", "generate a terraform file (main.tf)").Alias("g")
	commandGenerate.Arg("terraform-template", "path to a go template file").Required().ExistingFile()
	commandGenerate.Arg("context", "path to a yaml file").Required().ExistingFile()

	commandGenerateContext := app.Command("generate-context", "generate a context yaml file").Alias("c")
	commandGenerateContext.Flag("state", "(optional) path to a terraform.tfsate file").ExistingFiles()
	commandGenerateContext.Flag("template", "(optional) path to a go template file").ExistingFiles()
	commandGenerateContext.Flag("callback", "(optional) list of executable script file printing YAML to stdout").ExistingFile()
	commandGenerateContext.Arg("config-files", "(optional) list of yaml files").ExistingFiles()

	s := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch s {
	case (commandGenerate.FullCommand()):
		generate()
	case (commandGenerateContext.FullCommand()):
		generateContext()
	default:
		fmt.Println(app.Help)
	}
}

func generateContext() {
	fmt.Println("generate context")
}
func generate() {
	fmt.Println("generate")
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

func printVersion() string {
	return fmt.Sprintf("{\n"+
		"  version: \"%v\"\n"+
		"  commit-id: \"%v\"\n"+
		"  date: \"%v\"\n"+
		"}", version, commit, date)
}

func tfStringListFormater(a []interface{}) string {
	var result string = "[]"

	if a == nil {
		return result
	}
	if len(a) == 0 {
		return result
	}
	if len(a) == 1 {
		return fmt.Sprintf("[\"%v\"]", a[0])
	}

	result = "["
	for idx, val := range a {
		if idx < (len(a) - 1) {
			result = fmt.Sprintf("%v\"%v\", ", result, val)
		} else {
			result = fmt.Sprintf("%v\"%v\"]", result, val)
		}
	}

	return result
}
