package target

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Ensure we are running mstk in mstk root
*/
func EnsureMSTKRoot() {
	// Create bin directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	// Check that we are in MSTK
	base := filepath.Base(cwd)
	if base != "MSTK" && base != "mstk" {
		log.Panicf("Target is only allowed in root directory, but was ran inside %s", base)
	}
}

/*
Get all directories inside the root directory
*/
func GetDirectories(root string) ([]string, error) {
	var files []string

	f, err := os.Open(root)

	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

/*
Time the execution time of a function
*/
func TimeFn(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}

/*
Returns a list of services in mstk
*/
func GetMstkServicesList() []string {
	directories, err := GetDirectories("gomods/services/")
	if err != nil {
		log.Panic(err)
	}

	return directories
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Config for setupService function
*/
type setupServiceCfg struct {
	tag        string
	name       string
	dockerPath string
	runDir     string
}

/*
Apply kubectl command
*/
func applyKubectl(file, namespace string) {
	defer TimeFn(fmt.Sprintf("Deploy %s", file))()

	applyCmd := exec.Command(
		"kubectl", "apply", "-f", file, "-n", namespace,
	)
	log.Println("Applying to kubernetes")
	log.Printf("Running %s", applyCmd.String())

	err := applyCmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}
}

/*
Create a namespace in k8s
*/
func createNamespace(namespace string) {
	createCmd := exec.Command(
		"kubectl", "create", "namespace", namespace,
	)
	log.Println("Applying to kubernetes")
	log.Printf("Running %s", createCmd.String())

	err := createCmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}
}

/*
Delete a namespace in k8s
*/
func deleteNamespace(namespace string) {
	deleteCmd := exec.Command(
		"kubectl", "delete", "namespace", namespace,
	)
	log.Println("Applying to kubernetes")
	log.Printf("Running %s", deleteCmd.String())

	err := deleteCmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}
}

/*
Build go sources ready for docker
*/
func buildSourcesForDocker(targetFile, sourceDir string) {
	// TODO: GOOS, GOARCH options

	cmd := exec.Command("go", "build", "-o", targetFile, sourceDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=arm")

	log.Println("Running ", cmd.String())
	_, err := cmd.Output()
	if err != nil {
		log.Println(string((err.(*exec.ExitError).Stderr)))
		log.Panic(err.Error())
	}
}

/*
Write docker file template
*/
func writeDockerFile(path, service string) {
	var templateData struct {
		Time    string
		BinDir  string
		Package string
	}

	templateData.Package = service
	templateData.BinDir = path
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

	file, err := os.Create(fmt.Sprintf("%sdocker/Dockerfile.%s", path, service))
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
Setup a service
*/
func setupService(cfg setupServiceCfg) {
	defer TimeFn(fmt.Sprintf("Setting up '%s'", cfg.name))()

	buildCmd := exec.Command(
		"docker", "build",
		"--platform", "linux/arm64",
		"-t", fmt.Sprintf("%s%s:%s", cfg.tag, cfg.name, version),
		"-f", cfg.dockerPath,
		".",
	)
	inMinikube(buildCmd)
	buildCmd.Dir = cfg.runDir

	// TODO: Better error checking
	log.Printf("Building '%s'", cfg.name)
	log.Printf("Running %s", buildCmd.String())
	_, err := buildCmd.Output()
	// log.Println(string(out))
	if err != nil {
		log.Println(string((err.(*exec.ExitError).Stderr)))
		log.Panic(err)
	}

	pushCmd := exec.Command(
		"docker", "image", "push",
		fmt.Sprintf("%s:%s", cfg.name, version),
	)
	inMinikube(pushCmd)
	pushCmd.Dir = cfg.runDir

	log.Printf("Pushing '%s'", cfg.name)
	log.Printf("Running %s", buildCmd.String())
	_, err = pushCmd.Output()
	// log.Println(string(out))
	if err != nil {
		// Image push here fails sometimes with a repository error which isn't fatal, otherwise it is
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			fmt.Sprintf("An image does not exist locally with the tag: %s", cfg.name),
		) {
			log.Println(string((err.(*exec.ExitError).Stderr)))
			log.Panic(err)
		}
	}
}

/*
Adds environment values to a command
(Same as running eval $(minikube docker-env))
*/
func inMinikube(cmd *exec.Cmd) {
	minikubeEnv := exec.Command("minikube", "docker-env")

	out, err := minikubeEnv.Output()
	if err != nil {
		log.Println(err)
	}

	cmd.Env = os.Environ()

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "export") {
			// We got an export
			line, _ := strings.CutPrefix(scanner.Text(), "export ")

			split := strings.Split(line, "=")
			arg := strings.Trim(split[1], "\"")

			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", split[0], arg))
		}
	}
}

/*
Copies all files from one directory to another
*/
func copyDir(from, to string, ignoreExtensions []string) {
	f, err := os.Open(from)
	if err != nil {
		log.Panic(err)
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Copying directory '%s' to '%s'", from, to)

fileLoop:
	for _, file := range fileInfo {
		if file.IsDir() {
			// TODO: No symlinks etc.
			os.Mkdir(to+file.Name(), os.ModePerm)
			copyDir(from+file.Name(), to+file.Name(), ignoreExtensions)
		} else {
			for _, ext := range ignoreExtensions {
				if ext == filepath.Ext(file.Name()) {
					continue fileLoop
				}
			}

			log.Printf("\t%s", file.Name())

			// Read all content of src to data, may cause OOM for a large file.
			data, err := os.ReadFile(from + file.Name())
			if err != nil {
				log.Panic(err)
			}

			// Write data to dst
			err = os.WriteFile(to+file.Name(), data, fs.ModePerm)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

/*
Apply all mstk kubernetes files to the namespace
*/
func applyMstkK8s(projectName string) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating in %s", dirname)
	baseDir := dirname + "/mstk/k8s/"

	applyKubectl(baseDir, projectName)
}

/*
Executes a rolling restart for a namespace
*/
func restartKubernetes(namespace string) {
	defer TimeFn("Restart")()

	applyCmd := exec.Command(
		"kubectl", "rollout", "restart", "deployment", "-n", namespace,
	)
	log.Printf("Running %s", applyCmd.String())

	err := applyCmd.Run()
	if err != nil {
		log.Panic(err)
	}
}
