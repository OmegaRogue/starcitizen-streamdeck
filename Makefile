.PHONY: plugin generate build prepare
SHELL := /usr/bin/zsh
.ONESHELL:

plugin: build
	cp files/manifest.json out
	cp files/config.yaml out
	cp -r files/{Images,Sounds,PropertyInspector} out
	export version=$$(git describe --tags --abbrev=0 | sed 's/\([^-]*-g\)/r\1/;s/-/./g')
	sed -i "s/{{version}}/$${version}/g" out/manifest.json
	rm out/PropertyInspector/StarCitizen/*
	cd out
	rm starcitizen-streamdeck.streamdeckPlugin
	zip -r starcitizen-streamdeck.streamdeckPlugin *

prepare:
	mkdir -p out
generate: prepare
	go generate -v ./...
build: generate
	go build -o out/ starcitizen-streamdeck
lint:
	golangci-lint run