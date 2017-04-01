.DEFAULT_GOAL: all

.PHONY: all
all: mux

.PHONY: mux
mux:
	go build -o bin/mux examples/mux/main.go