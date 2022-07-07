// Copyright 2022. Motty Cohen
//
// Docker engine client library wrapper
//
package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/mottyc/yaaf-common/utils/collections"

	"github.com/docker/docker/api/types"
	docker_types "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

// General File utils
type DockerClient struct {
	dockerClient *client.Client
	ctx          context.Context
}

// New creates a docker engine client
func NewDockerClient() (*DockerClient, error) {
	if dc, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation()); err != nil {
		return nil, err
	} else {
		cli := &DockerClient{
			dockerClient: dc,
			ctx:          context.Background(),
		}
		return cli, nil
	}
}

// DefineContainer begins defining a docker container via a fluent interface.
func (c *DockerClient) CreateContainer(image string) *DockerContainer {
	return &DockerContainer{
		client:     c,
		image:      image,
		ports:      make(map[string]string),
		vars:       make(map[string]string),
		entryPoint: make([]string, 0),
		autoRemove: true,
	}
}

// FindContainerByName returns handle to an existing container from its name, or nil if the container was not found.
func (c *DockerClient) FindContainerByName(name string) (containerID string, err error) {

	list, err := c.dockerClient.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}

	// Search container by name
	for _, container := range list {
		if collections.Include(container.Names, "/"+name) {
			return container.ID, nil
		}
	}
	return "", nil
}

// createAndRunContainer creates a new docker container, initialize and runs it
func (c *DockerClient) createAndRunContainer(spec *DockerContainer) (containerID string, err error) {

	if containerID, err = c.createContainer(spec); err != nil {
		return
	}

	if err = c.dockerClient.ContainerStart(c.ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return
	} else {
		return containerID, nil
	}
}

// createContainer creates a new docker container based on the provided spec
func (c *DockerClient) createContainer(spec *DockerContainer) (string, error) {

	// verify image exists or pull it from the docker registry
	if _, err := c.verifyImage(spec.image); err != nil {
		return "", err
	}

	// Verify that the container does not exist
	if containerID, err := c.FindContainerByName(spec.name); err != nil {
		return "", err
	} else if len(containerID) > 0 {
		return containerID, fmt.Errorf("container %s is already running (%s)", spec.name, containerID)

	}

	// Set environment variables
	env := make([]string, 0)
	for k, v := range spec.vars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Set ports binding
	portBindings := nat.PortMap{}

	for k, v := range spec.ports {
		hostBinding := nat.PortBinding{HostIP: "0.0.0.0", HostPort: k}

		containerPort, err := nat.NewPort("tcp", v)
		if err != nil {
			return "", fmt.Errorf("unable to bind to port tcp:%s", v)
		}
		portBindings[containerPort] = []nat.PortBinding{hostBinding}
	}

	containerConfig := &docker_types.Config{
		Image: spec.image,
		Env:   env,
	}

	if len(spec.entryPoint) > 0 {
		containerConfig.Entrypoint = spec.entryPoint
	}

	hostConfig := &docker_types.HostConfig{
		AutoRemove:   spec.autoRemove,
		PortBindings: portBindings,
	}

	if container, err := c.dockerClient.ContainerCreate(c.ctx, containerConfig, hostConfig, nil, nil, spec.name); err != nil {
		return "", err
	} else {
		return container.ID, nil
	}
}

// verifyImage will verify that the docker image name exists or pull it from the docker repository
func (c *DockerClient) verifyImage(name string) (bool, error) {

	// Get list of all existing docker images
	images, err := c.dockerClient.ImageList(c.ctx, types.ImageListOptions{All: true})
	if err != nil {
		return false, err
	}

	//
	for _, image := range images {
		for _, imageName := range image.RepoTags {
			if name == imageName {
				return true, nil
			}
		}
	}

	// Image not exist in the local machine, need to pull it from the docker registry
	reader, err := c.dockerClient.ImagePull(c.ctx, name, types.ImagePullOptions{})
	if err != nil {
		return false, err
	}

	defer func() {
		_ = reader.Close()
	}()

	if _, err = io.Copy(os.Stdout, reader); err != nil {
		return false, fmt.Errorf("error pulling image: %s", err)
	}
	return false, nil
}

// RemoveContainer stop, kill and remove the container
func (c *DockerClient) RemoveContainer(containerID string) error {
	return c.dockerClient.ContainerRemove(c.ctx, containerID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
}

// RemoveContainer stop, kill and remove the container
func (c *DockerClient) GetContainerState(containerID string) (string, error) {

	list, err := c.dockerClient.ContainerList(c.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}

	// Search container by id
	for _, item := range list {
		if item.ID == containerID {
			return item.State, nil
		}
	}

	return "", fmt.Errorf("container %s not found", containerID)
}
