// Copyright 2022. Motty Cohen
//
// Docker engine client library tests
//
package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDockerClientRunContainer(t *testing.T) {

	// Create client
	cli, err := NewDockerClient()
	assert.Nil(t, err)

	// Create container and run it
	id, err := cli.CreateContainer("busybox:latest").
		Name("busybox").
		EntryPoint("tail", "-f", "/dev/null").
		Label("environment", "test").
		Label("group", "core").
		Run()

	// Check state
	state, er := cli.GetContainerState(id)
	assert.Nil(t, er)
	assert.Equal(t, state, "running")

	// Find by name
	containerId, err := cli.FindContainerByName("busybox")
	assert.Nil(t, err)
	assert.Equal(t, containerId, id)

	// List by label
	list, err := cli.ListContainersByLabel("group", "core")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))

	err = cli.RemoveContainer(containerId)
	assert.Nil(t, err)

}
