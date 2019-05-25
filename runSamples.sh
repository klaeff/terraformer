mkdir -p gen/examples

go run terraformer.go ctx > gen/examples/context1.yml
go run terraformer.go ctx --state=./examples/context/terraform.tfstate > gen/examples/context2.yml
go run terraformer.go ctx --callback=./examples/context/callback.sh > gen/examples/context3.yml
go run terraformer.go ctx ./examples/context/config1.yml > gen/examples/context4.yml

go run terraformer.go gen ./examples/aws/tf.template ./examples/aws/context.yml > gen/examples/tf.main

error=$(cat gen/examples/* | grep error)

if [ "$error" == "" ]; then
  ls -al gen/examples
else
  echo "--$error--"
fi
