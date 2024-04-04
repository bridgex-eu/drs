package dockerhelper

import (
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/remote"
	"github.com/schollz/progressbar/v3"
)

func isSuperUser(client *command.MachineClient) bool {
	cmdStr := "id -un | grep -qx 'root' || command -v sudo || command -v su || exit 1"

	return client.Command(cmdStr).Run() == nil
}

func logCmd(cmd *remote.Cmd) {
	slog.Info(
		fmt.Sprintf("Running command on %s", cmd.Addr().String()),
		slog.String("cmd", cmd.String()),
	)
}

func isDockerInstalled(client *command.MachineClient) bool {
	return client.Command("docker", "-v").Run() == nil
}

func ExecuteDockerCmd(client *command.MachineClient, stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	cmd := client.Command("docker", args...)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func installDocker(client *command.MachineClient) error {
	installCmd := `
        curl -fsSL https://get.docker.com | sudo sh || \
        wget -qO- https://get.docker.com | sudo sh || \
        exit 1
    `

	if err := client.Command(installCmd).Run(); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	setupCmd := `
        sudo systemctl enable docker.service && \
        sudo systemctl enable containerd.service && \
        docker run --privileged --rm tonistiigi/binfmt --install all
    `

	if err := client.Command(setupCmd).Run(); err != nil {
		return fmt.Errorf("failed to setup Docker services and binfmt: %w", err)
	}

	return nil
}

var drsNetwork = "drs"

var containerFlags = []string{
	"-d",
	"--log-opt max-size=10m",
	"--log-opt max-file=10",
	"--restart=unless-stopped",
	"--network " + drsNetwork,
}

func RemoveUnusedImages(
	client *command.MachineClient,
) error {
	cmd := client.Command("docker", "image", "prune", "-a")
	logCmd(cmd)

	return cmd.Run()
}

func RenameContainer(
	client *command.MachineClient,
	name, newName string,
) error {
	cmd := client.Command("docker", "rename", name, newName)
	logCmd(cmd)

	return cmd.Run()
}

func StartContainer(
	client *command.MachineClient,
	name string,
) error {
	cmd := client.Command("docker", "start", name)
	logCmd(cmd)

	return cmd.Run()
}

func StopContainer(
	client *command.MachineClient,
	name string,
) error {
	cmd := client.Command("docker", "stop", name)
	logCmd(cmd)

	return cmd.Run()
}

func RemoveContainer(
	client *command.MachineClient,
	name string,
	force bool,
) error {
	args := []string{"rm", "-v"}
	if force {
		args = append(args, "-f")
	}

	args = append(args, name)

	if err := StopContainer(client, name); err != nil {
		return err
	}

	cmd := client.Command("docker", args...)
	logCmd(cmd)

	return cmd.Run()
}

func RunContainer(
	client *command.MachineClient,
	image, name string,
	ports, volumes, labels, envs []string,
	args ...string,
) error {
	cmdArgs := []string{"run"}

	cmdArgs = append(cmdArgs, "--name", name, "--network-alias", name)
	cmdArgs = append(cmdArgs, containerFlags...)

	for _, label := range labels {
		cmdArgs = append(cmdArgs, "-l", "'"+label+"'")
	}

	for _, port := range ports {
		cmdArgs = append(cmdArgs, "-p", port)
	}

	for _, volume := range volumes {
		cmdArgs = append(cmdArgs, "-v", volume)
	}

	for _, env := range envs {
		cmdArgs = append(cmdArgs, "-e", env)
	}

	cmdArgs = append(cmdArgs, image)
	cmdArgs = append(cmdArgs, args...)

	cmd := client.Command("docker", cmdArgs...)
	logCmd(cmd)

	return cmd.Run()
}

func CreateDockerNetwork(client *command.MachineClient, name string) error {
	cmd := client.Command("docker", "network", "create", name)
	logCmd(cmd)

	return cmd.Run()
}

func IsContainerExist(client *command.MachineClient, name string) bool {
	nameFormat := fmt.Sprintf("name=^%s$", name)
	output, err := client.Command("docker", "ps", "-a", "-q", "-f", nameFormat).Output()

	// Check if the command execution was successful and if any output was generated
	// A non-empty output indicates the container is running
	return err == nil && output != ""
}

func IsNetworkExist(client *command.MachineClient, name string) bool {
	nameFormat := fmt.Sprintf("name=^%s$", name)
	output, err := client.Command("docker", "network", "ls", "-q", "-f", nameFormat).CombinedOutput()

	// Check if the command execution was successful and if any output was generated
	// A non-empty output indicates the network exists
	return err == nil && output != ""
}

func LocalImageSize(name string) (int, error) {
	out, err := exec.Command("docker", "inspect", "-f", "{{ .Size }}", name).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("Failed to get docker image size: %w. Details: %s", err, out)
	}

	return strconv.Atoi(strings.TrimSpace(string(out)))
}

func SendDockerImage(client *command.MachineClient, image string) error {
	cmd := exec.Command("docker", "save", image)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to get stdout pipe for docker save: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start docker save: %w", err)
	}

	size, err := LocalImageSize(image)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(
		int64(size),
		"Sending",
	)

	reader, writer := io.Pipe()

	remoteCmd := client.Command("docker", "load")
	remoteCmd.Stdin = reader

	// Use a goroutine for copying the output to the multiWriter
	go func() {
		multiWriter := io.MultiWriter(bar, writer)
		if _, err := io.Copy(multiWriter, output); err != nil {
			fmt.Println("Failed to copy docker image to multiWriter:", err)
		}
		writer.Close() // Close the writer to signal the remote command that the input is done
	}()

	if err := remoteCmd.Run(); err != nil {
		return fmt.Errorf("Failed to execute remote docker load: %w", err)
	}

	return nil
}
