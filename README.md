# terraformer
A go program that generates terraform files using go templates

![terraformer](doc/terraformer-planet.jpg)

## installation (osx)

```
brew install sklevenz/skl/terraformer 
```

## usage

```
usage: terraformer [<flags>] <command> [<args> ...]

A go program that generates terraform files using go templates

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  generate <terraform-template> <context>
    generate a terraform file (main.tf)

  generate-context [<flags>] [<config-files>...]
    generate a context yaml file
```


## motivation

Terraform has its own proprietary programming and template syntax. Many infrastructure projects start small and often grow after a while. Then the complexity of a proprietary technology could become an issue. 

The thesis is that "standard" programming technologies can handle this complexity better. Example technologies are Ruby, Python, node.js ... Standard means there are simply more developers familiar with this technologies. That's it.

This solution uses go as technology and the go templates. Go, because of the single binary is simple to use. 

## concept





The idea is:

- step 1
  - collect all configuration within a single context yaml file
  - configuration can come from
    - environment variables
    - config files
    - terraform state files
    - command executions
    - ...  
- step 2
  - use a simple go template
  - and generate a main.tf file using the context data
- step 3
  - apply the main.tf file to terraform 

The template file should be simple. Ideally it is a flat list of resources. Programming logic and variable substitution is done by go templating.
The context is a single source of data. Basically it is a yaml data structure which is accessible in the go template as `{{.}}`. This context is maybe generated as well and get's data from all kind of sources. Static configuration, command executions, environment variables, other yaml files, terraform state files, etc. The result should be a main.tf containing all resources and all data. Apply this with terraform.

![terraformer](doc/terraformer.png)

## best practices

Data sources, variable, loops, dependencies, modules ... all of this makes maintenance of terraform difficult. We recommend the following

- create more but smaller templates 
- deploy in chunks and have a defined order
- practice TDI (test driven infrastructure)
  - create a test account
  - deploy form scratch, update and delete

## examples

terraform template file

```
provider "aws" {
  access_key = "{{ .context.access_key }}"
  secret_key = "{{ .context.secret_key }}"
  region     = "us-west-1"
}
```

context yaml file

```
context:
  access_key: 123
  secret_key: "abc"
```

## features

Go template processing is can be easily extended by special functions. The first command is `tfStringListFormater`. This was required because of kubectl of kubernetes returns list of IPs without quotes which are not accepted by terraform.

| feature | description | example |
|---------|-------------|---------|
| tfStringListFormater | formats a list with quoted elements | [1,2,3] -> ["1","2","3"] |
| more to come | provide a pull request | f(x)  |

## try out 

### context generation

```
go run terraformer.go generate-context
go run terraformer.go generate-context --state=./examples/context/terraform.tfstate
go run terraformer.go generate-context --callback=./examples/context/callback.sh
go run terraformer.go generate-context ./examples/context/config1.yml
```

Combine all together.

### tf main generation

```
go run terraformer.go gen ./examples/aws/tf.template ./examples/aws/config.yml
```

## build

```
go test
go build
```

## todo

- (ok) introduce flags, kingpin, etc for better command line experience
  - (ok) unit test
- implement context generation
  - (ok) basics, test
  - (ok) config, test
  - (ok) state, test
  - callback, test 
  - template, test
- (ok) implement tf generation
  - (ok) basics, test
- more terraform samples
- bash completion support (via homebrew)
- man pages support

![terraformer](doc/terraformer-logo-small.png)

