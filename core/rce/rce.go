package rce

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Language string

const (
	PYTHON Language = "python"
	JAVA   Language = "java"
	C      Language = "c"
	CPP    Language = "cpp"
)

const (
	MemoryLimit = 10000000 // 10MB (Minimum memory limit allowed is 6MB by Docker)
	CPUQuota    = 100000   // 0.1 CPU
)

// RunProgram runs the given program in a Docker container with the specified language.
func RunProgram(program string, language Language) (string, string, error) {
	containerImage, cmd, containerName, err := getContainerConfig(program, language)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()

	apiClient, err := createNewAPIClient()
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return "", "", err
	}

	resp, err := createContainer(ctx, apiClient, containerImage, cmd, containerName)
	if err != nil {
		log.Printf("Failed to create Docker container: %v", err)
		return "", "", err
	}

	if err := startContainer(ctx, apiClient, resp.ID); err != nil {
		log.Printf("Failed to start Docker container: %v", err)
		return "", "", err
	}
	stdout, stderr, err := getContainerLogs(ctx, apiClient, resp.ID)
	if err != nil {
		log.Printf("Failed to get Docker container logs: %v", err)
		return "", "", err
	}

	if err := removeContainer(ctx, apiClient, resp.ID); err != nil {
		log.Printf("Failed to remove Docker container: %v", err)
		return "", "", err
	}

	return stdout, stderr, err
}

// getContainerConfig returns the Docker container image, command, and name based on the programming language.
func getContainerConfig(program string, language Language) (string, []string, string, error) {
	containerImage := ""
	cmd := []string{}
	containerName := ""

	switch language {
	case PYTHON:
		containerImage = "python"
		cmd = []string{"python", "-c", program}
		containerName = "python-code-runner"
	case JAVA:
		containerImage = "openjdk"
		cmd = []string{"sh", "-c", "echo '" + program + "' > Main.java && javac Main.java && java Main"}
		containerName = "java-code-runner"
	case C:
		containerImage = "gcc"
		cmd = []string{"sh", "-c", "echo '" + program + "' > main.c && gcc main.c -o main && ./main"}
		containerName = "c-code-runner"
	case CPP:
		containerImage = "gcc"
		cmd = []string{"sh", "-c", "echo '" + program + "' > main.cpp && g++ main.cpp -o main && ./main"}
		containerName = "cpp-code-runner"
	default:
		log.Printf("Unsupported language: %s", language)
		return "", cmd, "", fmt.Errorf("unsupported language: %s", language)
	}

	return containerImage, cmd, containerName, nil
}

// createNewAPIClient creates a new Docker API client.
func createNewAPIClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// createContainer creates a new Docker container with the specified image, command, and name.
func createContainer(ctx context.Context, apiClient *client.Client, containerImage string, cmd []string, containerName string) (container.CreateResponse, error) {
	pidsLimit := int64(100)

	return apiClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:           containerImage,
			Cmd:             cmd,
			AttachStdout:    true,
			AttachStderr:    true,
			NetworkDisabled: true,
			User:            "nobody", // Run as non-root user
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    MemoryLimit,
				CPUQuota:  CPUQuota,
				PidsLimit: &pidsLimit,
			},
			NetworkMode:    "none", // Disable networking
			ReadonlyRootfs: true,   // Make filesystem read-only
			SecurityOpt: []string{
				"no-new-privileges",  // Prevent escalation of privileges
			},
		},
		nil,
		nil,
		containerName,
	)
}

// startContainer starts the Docker container with the specified ID.
func startContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	if err := apiClient.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return err
	}

	// Wait for the Docker container to finish running
	statusCh, errCh := apiClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

// getContainerLogs returns the logs from the Docker container with the specified ID.
func getContainerLogs(ctx context.Context, apiClient *client.Client, containerID string) (string, string, error) {
	out, err := apiClient.ContainerLogs(
		ctx,
		containerID,
		container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     true,
			Details:    true,
		})
	if err != nil {
		return "", "", err
	}
	defer out.Close()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	_, err = stdcopy.StdCopy(stdout, stderr, out)
	if err != nil {
		return "", "", err
	}

	// Combine stdout and stderr logs
	return stdout.String(), stderr.String(), nil
}

// removeContainer removes the Docker container with the specified ID.
func removeContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	return apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
}
