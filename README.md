# yaaf-common-docker

[![Build](https://github.com/go-yaaf/yaaf-common-docker/actions/workflows/build.yml/badge.svg)](https://github.com/go-yaaf/yaaf-common-docker/actions/workflows/build.yml)

A wrapper library of the Go client for the Docker Engine API.

## About
This library is used to simplify docker container orchestration with Docker engine (using fluent API) and it is mainly used for integration tests.

This library is built around the [official Go SDK for Docker](https://github.com/docker/go-docker)
and includes a subset of docker CLI functions to be used for integration tests.

#### Adding dependency

```bash
$ go get -v -t github.com/go-yaaf/yaaf-common-docker ./...
```

### Library usage
```go
package main

import (
	"github.com/go-yaaf/yaaf-common-docker/docker"
	"log"
)

func main() {

	// Create client
	cli, err := docker.NewDockerClient()
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Create container and run it
	id, err := cli.CreateContainer("busybox:latest").
		Name("busybox").
		EntryPoint("/bin/echo", "busybox", "foo").
		Label("environment", "test").
		Label("group", "core").
		Run()

	// Check state
	if state, err := cli.GetContainerState(id); err != nil {
		log.Fatal(err)
    } else {
		log.Println(state)
    }

	// Find by name
	if containerId, err := cli.FindContainerByName("busybox"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("found", containerId)
	}

	// List by label
	if list, err := cli.ListContainersByLabel("group", "core"); err != nil {
		log.Fatal(err)
	} else {
		for i, item := range list {
			log.Println(i, item)
        }
	}

	// Remove the container
	if err := cli.RemoveContainer(id); err != nil {
		log.Fatal(err)
	} else {
		log.Println("container", id, "removed")
	}
}
```