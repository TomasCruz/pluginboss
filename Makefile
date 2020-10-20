all: generate_grpc_code build build_plugins

clean:
	rm -f bin/*
	go clean

generate_grpc_code:
	protoc --go_out=plugins=grpc:plugininfo plugininfo/converter.proto

build:
	go build -o bin/plugin_boss

build_plugins:
	go build -o bin/plugin_one ./plugin-one
	go build -o bin/plugin_two ./plugin-two
	go build -o bin/plugin_three ./plugin-three
	go build -o bin/plugin_four ./plugin-four
