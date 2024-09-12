.PHONY: plugin generate build prepare
SHELL := /usr/bin/zsh
.ONESHELL:

plugin: build
	cp files/manifest.json out/codes.omegavoid.starcitizen.sdPlugin
	cp files/config.yaml out/codes.omegavoid.starcitizen.sdPlugin
	cp -r files/{Images,Sounds,PropertyInspector,*.gohtml} out/codes.omegavoid.starcitizen.sdPlugin
	export version=$$(git describe --tags --abbrev=0 | sed 's/\([^-]*-g\)/r\1/;s/-/./g')
	sed -i "s/{{version}}/$${version}/g" out/codes.omegavoid.starcitizen.sdPlugin/manifest.json
	rm out/codes.omegavoid.starcitizen.sdPlugin/PropertyInspector/StarCitizen/*
	cd out
	rm starcitizen-streamdeck.streamdeckPlugin
	zip -r starcitizen-streamdeck.streamdeckPlugin codes.omegavoid.starcitizen.sdPlugin

prepare:
	mkdir -p out/codes.omegavoid.starcitizen.sdPlugin
generate: prepare
	go generate -v ./...
build: generate
	go build -o out/ starcitizen-streamdeck
lint:
	golangci-lint run