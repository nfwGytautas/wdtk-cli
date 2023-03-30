package target

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Flags for setup target
*/
var SetupFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "platform",
		Usage: "platform for docker",
	},
}

/*
Execute setup target
*/
func SetupAction(ctx *cli.Context) {
	defer TimeFn("Setup")()
	// TODO: Find the MSTK installation path automatically

	log.Println("Running setup")
	EnsureMSTKRoot()

	log.Println("Creating bin")
	os.Mkdir("bin/", os.ModePerm)

	log.Println("Compiling services")
	services := GetMstkServicesList()

	log.Printf("Found %v services %v", len(services), services)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		go compileService(service, &wg)
	}
	wg.Wait()

	log.Println("Setup done, your minikube environment should have mstk microservices up and running")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Version for images
*/
const version = "0.0.0"

/*
Compiles a single service
*/
func compileService(path string, wg *sync.WaitGroup) {
	defer TimeFn(fmt.Sprintf("Preparing '%s'", path))()
	defer wg.Done()

	log.Printf("Compiling %s", path)

	serviceName := filepath.Base(path)
	targetDir := fmt.Sprintf("./bin/%s", serviceName)
	sourceDir := fmt.Sprintf("./gomods/%s/", serviceName)

	// TODO: GOOS, GOARCH options
	cmd := exec.Command("go", "build", "-o", targetDir, sourceDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=arm")

	log.Println("Running ", cmd.String())
	_, err := cmd.Output()
	if err != nil {
		log.Println(string((err.(*exec.ExitError).Stderr)))
		log.Panic(err.Error())
	}

	// Generate docker files
	writeDockerFile(serviceName)

	// Push to minikube
	setupService(serviceName)

	// Apply kubectl commands
	applyKubectl(serviceName)
}

/*
Write docker file template
*/
func writeDockerFile(service string) {
	var templateData struct {
		Time    string
		Package string
	}

	templateData.Package = service
	templateData.Time = time.Now().String()

	template, err := template.New("dockerfile").Parse(`
# Automatically generated docker file
# Generated at {{.Time}}

# Go as minimal as possible
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY ./bin/{{.Package}} /{{.Package}}

EXPOSE 8080

USER root

ENTRYPOINT ["/{{.Package}}"]
`)

	if err != nil {
		log.Panic(err)
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(fmt.Sprintf("./bin/Dockerfile.%s", service))
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Panic(err)
	}
	file.Sync()
}

/*
Adds environment values to a command
(Same as running eval $(minikube docker-env))
*/
func inMinikube(cmd *exec.Cmd) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_TLS_VERIFY=1")
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:56432")
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CERT_PATH=%s/.minikube/certs", dirname))
	cmd.Env = append(cmd.Env, "MINIKUBE_ACTIVE_DOCKERD=\"minikube\"")
}

/*
Setup a service
*/
func setupService(service string) {
	defer TimeFn(fmt.Sprintf("Setting up '%s'", service))()

	buildCmd := exec.Command(
		"docker", "build",
		"--platform", "linux/arm64",
		"-t", fmt.Sprintf("mstk/%s:%s", service, version),
		"-f", fmt.Sprintf("./bin/Dockerfile.%s", service),
		".",
	)
	inMinikube(buildCmd)

	// TODO: Better error checking
	log.Printf("Building '%s'", service)
	log.Printf("Running %s", buildCmd.String())
	_, err := buildCmd.Output()
	// log.Println(string(out))
	if err != nil {
		log.Println(string((err.(*exec.ExitError).Stderr)))
		log.Panic(err)
	}

	pushCmd := exec.Command(
		"docker", "image", "push",
		fmt.Sprintf("%s:%s", service, version),
	)
	inMinikube(pushCmd)

	log.Printf("Pushing '%s'", service)
	log.Printf("Running %s", buildCmd.String())
	_, err = pushCmd.Output()
	// log.Println(string(out))
	if err != nil {
		// Image push here fails sometimes with a repository error which isn't fatal, otherwise it is
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			fmt.Sprintf("An image does not exist locally with the tag: %s", service),
		) {
			log.Println(string((err.(*exec.ExitError).Stderr)))
			log.Panic(err)
		}
	}
}

/*
Apply kubectl commands for a service
*/
func applyKubectl(service string) {
	defer TimeFn(fmt.Sprintf("Applying %s", service))()

	applyCmd := exec.Command(
		"kubectl", "apply", "-f", fmt.Sprintf("kubes/%s/", service),
	)
	log.Println("Applying to kubernetes")
	log.Printf("Running %s", applyCmd.String())

	err := applyCmd.Run()
	if err != nil {
		log.Panic(err)
	}
}
