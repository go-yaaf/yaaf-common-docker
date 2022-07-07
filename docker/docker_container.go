// Copyright 2022. Motty Cohen
//
// Docker container specification
//
package docker

// DockerContainer is used to construct docker container spec for the docker engine.
type DockerContainer struct {
	client     *DockerClient     // Docker client
	image      string            // Docker image
	name       string            // Container name
	ports      map[string]string // Container ports mapping
	vars       map[string]string // Environment variables
	entryPoint []string          // Entry point
	autoRemove bool              // Automatically remove container when stopped (default: true)
}

// Name sets the container name.
func (c *DockerContainer) Name(value string) *DockerContainer {
	c.name = value
	return c
}

// Port adds a port mapping
func (c *DockerContainer) Port(external, internal string) *DockerContainer {
	c.ports[external] = internal
	return c
}

// Ports adds multiple port mappings
func (c *DockerContainer) Ports(ports map[string]string) *DockerContainer {
	for k, v := range ports {
		c.ports[k] = v
	}
	return c
}

// Var adds an environment variable
func (c *DockerContainer) Var(key, value string) *DockerContainer {
	c.vars[key] = value
	return c
}

// Vars adds multiple environment variables
func (c *DockerContainer) Vars(vars map[string]string) *DockerContainer {
	for k, v := range vars {
		c.vars[k] = v
	}
	return c
}

// EntryPoint sets the entrypoint arguments of the container.
func (c *DockerContainer) EntryPoint(args ...string) *DockerContainer {
	c.entryPoint = append(c.entryPoint, args...)
	return c
}

// AutoRemove determines whether to automatically remove the container when it has stopped
func (c *DockerContainer) AutoRemove(value bool) *DockerContainer {
	c.autoRemove = value
	return c
}

// Run creates and runs the container, returning the container ID.
func (c *DockerContainer) Run() (containerID string, err error) {
	return c.client.createAndRunContainer(c)
}
