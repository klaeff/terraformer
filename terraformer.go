package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/alecthomas/kingpin/v2"
	"gopkg.in/yaml.v3"
)

var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

var (
	app = kingpin.New("terraformer", "A go program that generates terraform files using go templates")

	commandGenerate = app.Command("generate", "generate a terraform file (main.tf), alias=gen").Alias("gen")
	templateFile    = commandGenerate.Arg("terraform-template", "path to a go template file").Required().ExistingFile()
	contextFile     = commandGenerate.Arg("context", "path to a yaml file").Required().ExistingFile()

	commandGenerateContext = app.Command("generate-context", "generate a context yaml file, alias=ctx").Alias("ctx")
	stateFlag              = commandGenerateContext.Flag("state", "(optional) path to a terraform.tfsate file").Short('s').ExistingFile()
	templateFlag           = commandGenerateContext.Flag("template", "(optional) path to a go template file").Short('t').ExistingFile()
	callbackFlag           = commandGenerateContext.Flag("callback", "(optional) list of executable script file printing YAML to stdout").Short('c').ExistingFile()
	configFiles            = commandGenerateContext.Arg("config-files", "(optional) list of yaml files").ExistingFiles()
)

func main() {
	app.Version(printVersion()).Author("Stephan Klevenz")
	app.HelpFlag.Short('h')

	s := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch s {
	case (commandGenerate.FullCommand()):
		generate(*templateFile, *contextFile)
	case (commandGenerateContext.FullCommand()):
		generateContext(*stateFlag, *templateFlag, *callbackFlag, *configFiles)
	default:
		fmt.Println(app.Help)
	}

}

func generate(templateFile string, contextFile string) {
	funcMap := template.FuncMap{
		"tfStringListFormater": tfStringListFormater,
		"tfCallback":           tfCallback,
	}

	templateName := path.Base(templateFile)
	template, err := template.New(templateName).Funcs(funcMap).ParseFiles(templateFile)

	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(-1)
	}

	data, err := ioutil.ReadFile(contextFile)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(-1)
	}

	var context interface{}
	err = yaml.Unmarshal(data, &context)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(-1)
	}

	err = template.Execute(os.Stdout, context)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(-1)
	}
}

func parseEnvironment() map[string]string {
	parseEnvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}

	return parseEnvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})
}

func readJsonFile(filePath string) interface{} {
	var jsonData interface{}

	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("error %v\n ", err)
		os.Exit(-1)
	}
	err = json.Unmarshal(jsonFile, &jsonData)
	if err != nil {
		log.Fatalf("error %v\n", err)
		os.Exit(-1)
	}

	return jsonData
}

func readYamlFile(filePath string) interface{} {
	var yamlData interface{}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("error %v\n ", err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		log.Fatalf("error %v\n", err)
		os.Exit(-1)
	}

	return yamlData
}

func readYamlCallback(scriptFile string) interface{} {
	var err error
	yamlBytes, err := exec.Command(scriptFile).Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	var yamlData interface{}

	err = yaml.Unmarshal([]byte(yamlBytes), &yamlData)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(-1)
	}

	return yamlData
}

func generateContext(stateFlag string, templateFlag string, callbackFlag string, configFiles []string) {
	topNodes := make(map[string]interface{})
	topNodes["env"] = parseEnvironment()

	if stateFlag != "" {
		topNodes["state"] = readJsonFile(stateFlag)
	}

	if templateFlag != "" {
		topNodes["templates"] = "not implemented"
	}

	if callbackFlag != "" {
		topNodes["callback"] = readYamlCallback(callbackFlag)
	}

	for idx, val := range configFiles {
		cfg := readYamlFile(val)
		x := fmt.Sprintf("config%v", idx)
		topNodes[x] = cfg
	}

	ctx := make(map[string]interface{})
	ctx["context"] = topNodes

	yamlBytes, err := yaml.Marshal(&ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(string(yamlBytes))
}

func printVersion() string {
	return fmt.Sprintf("{\n"+
		"  version: \"%v\"\n"+
		"  commit-id: \"%v\"\n"+
		"  date: \"%v\"\n"+
		"}", version, commit, date)
}

func tfCallback(scriptFile string) string {

	// call script relative to template file
	scriptFile = path.Join(path.Dir(*templateFile), scriptFile)

	out, err := exec.Command(scriptFile).Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return strings.TrimSuffix(string(out), "\n")
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
