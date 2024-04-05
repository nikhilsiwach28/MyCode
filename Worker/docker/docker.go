package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient interface {
	RunContainer(code, language string) (string, error)
}

type Client struct {
	cli *client.Client
}


func NewDockerClient() (DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) RunContainer(code, language string) (string, error) {
	// Prepare language-specific Docker container configuration and command
	var image string
	var cmd []string

	switch language {
	case "python":
		image = "python:latest"
		cmd = []string{"python", "-c", code}
	case "cpp":
		image = "gcc:latest"
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > main.cpp && gcc -o main main.cpp && ./main", code)}
	case "java":
		image = "openjdk:latest"
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > Main.java && javac Main.java && java Main", code)}
	case "rust":
		image = "rust:latest"
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > main.rs && rustc -o main main.rs && ./main", code)}
	case "golang":
		image = "golang:latest"
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > main.go && go run main.go", code)}
	default:
		return "", fmt.Errorf("Unsupported language", )
	}

	// Create a new container
	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start the container
	if err := c.cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Wait for the container to finish
	statusCh, errCh := c.cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	// Retrieve container logs
	out, err := c.cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy logs to a string
	logs, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	// Remove the container
	if err := c.cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{}); err != nil {
		return "", err
	}

	return string(logs), nil
}
