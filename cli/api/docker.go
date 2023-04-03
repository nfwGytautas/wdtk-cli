package api

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nfwGytautas/mstk/cli/common"
)

// PUBLIC TYPES
// ========================================================================

/*
Docker object
*/
type Docker struct {
	context string
	tag     string
}

// PRIVATE TYPES
// ========================================================================

/*
Dockerfile template
*/
const dockerfileTemplate = `
# Automatically generated docker file
# Generated at {{.Time}}

# Go as minimal as possible
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY ./{{.BinDir}}{{.Package}} /{{.Package}}

EXPOSE 8080

USER root

ENTRYPOINT ["/{{.Package}}"]
`

// PUBLIC FUNCTIONS
// ========================================================================

/*
Create docker object
*/
func CreateDocker(contextDir, tag string) Docker {
	return Docker{context: contextDir, tag: tag}
}

/*
Write a dockerfile template

context argument is the root directory
*/
func (d *Docker) WriteTemplate(service, binDir string) error {
	var templateData struct {
		BinDir  string
		Time    string
		Package string
	}

	templateData.BinDir = binDir
	templateData.Package = service
	templateData.Time = time.Now().String()

	template, err := template.New("dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	if err != nil {
		return err
	}

	dockerRoot := fmt.Sprintf("%s/docker/", d.context)
	err = os.MkdirAll(dockerRoot, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%sDockerfile.%s", dockerRoot, service))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

/*
Build and push a docker image
*/
func (d *Docker) BuildAndPush(service string) error {
	defer common.TimeCurrentFn()

	imageName := fmt.Sprintf("%s:%s", service, common.DockerVersion)

	buildCmd := exec.Command(
		"docker", "build",
		"--platform", "linux/arm64",
		"-t", fmt.Sprintf("%s/%s", d.tag, imageName),
		"-f", fmt.Sprintf("%s/docker/Dockerfile.%s", d.context, service),
		".",
	)
	pushCmd := exec.Command(
		"docker", "image", "push", imageName,
	)

	// Setup commands
	err := inMinikube([]*exec.Cmd{buildCmd, pushCmd})
	if err != nil {
		return err
	}

	buildCmd.Dir = d.context
	pushCmd.Dir = d.context

	// Execute commands
	err = common.ExecCmd(buildCmd)
	if err != nil {
		return err
	}

	err = common.ExecCmd(pushCmd)
	if err != nil {
		// Image push here fails sometimes with a repository error which isn't fatal, otherwise it is
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			fmt.Sprintf("An image does not exist locally with the tag: %s", service),
		) {
			return err
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================

/*
Adds environment values to a command
(Same as running eval $(minikube docker-env))
*/
func inMinikube(commands []*exec.Cmd) error {
	minikubeEnv := exec.Command("minikube", "docker-env")

	out, err := minikubeEnv.Output()
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		cmd.Env = os.Environ()
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "export") {
			// We got an export
			line, _ := strings.CutPrefix(scanner.Text(), "export ")

			split := strings.Split(line, "=")
			arg := strings.Trim(split[1], "\"")

			for _, cmd := range commands {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", split[0], arg))
			}
		}
	}

	return nil
}
