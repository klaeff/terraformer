# terraformer
A go program that generates terraform using go templates

## usage

terraformer terraform.go.template config.yml 

The config.yaml will be read into a hashmap and can then be processed within the go template file.

I don't like terraforms programming syntax and think using a generic template language is the better choice. The output will be a simple, flat and  reliable terraform file which can be applyed by terraform without getting side effects. 

Why not ruby or phyton which has nice templating as well? Because of you don't get it as a single binary.


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

## try out 
```
go run terraformer.go ./examples/tf.template ./examples/conterxt.yml
```

## build

go test
go build


## installation (osx)

```
brew install sklevenz/skl/terraformer 
```