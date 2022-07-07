// Copyright 2022. Motty Cohen
//
// Docker engine client library tests
//
package test

import (
	"testing"

	"github.com/mottyc/yaaf-common-docker/docker"
	"github.com/stretchr/testify/assert"
)

func TestRunContainer(t *testing.T) {

	// Create engine
	de, err := docker.NewDockerClient()
	assert.Nil(t, err)

	// Create container and run it
	id, err := de.CreateContainer("busybox:latest").
		Name("busybox").
		EntryPoint("tail", "-f", "/dev/null").
		Run()

	state, er := de.GetContainerState(id)
	assert.Nil(t, er)
	assert.Equal(t, state, "running")

	containerId, err := de.FindContainerByName("busybox")
	assert.Nil(t, err)
	assert.Equal(t, containerId, id)

	err = de.RemoveContainer(containerId)
	assert.Nil(t, err)

}
