package api

import (
	"os/exec"
	"strings"

	"github.com/nfwGytautas/mstk/cli/common"
)

// PUBLIC TYPES
// ========================================================================

/*
Kubernetes object
*/
type Kubernetes struct {
	namespace string
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Create K8S object
*/
func CreateK8s(namespace string) Kubernetes {
	return Kubernetes{namespace: namespace}
}

/*
Apply kubectl command
*/
func (k *Kubernetes) Apply(file string) error {
	defer common.TimeCurrentFn()()

	cmd := exec.Command(
		"kubectl", "apply", "-f", file, "-n", k.namespace,
	)
	return common.ExecCmd(cmd)
}

/*
Delete kubectl command
*/
func (k *Kubernetes) Delete(file string) error {
	defer common.TimeCurrentFn()()

	cmd := exec.Command(
		"kubectl", "delete", "-f", file, "-n", k.namespace,
	)
	err := common.ExecCmd(cmd)

	if err != nil {
		// Not found is not an actual error, just it doesn't exist which is fine since we are cleaning up anyway
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			"not found",
		) {
			return err
		}
	}

	return nil
}

/*
Create namespace kubectl command
*/
func (k *Kubernetes) CreateNamespace() error {
	defer common.TimeCurrentFn()()

	cmd := exec.Command(
		"kubectl", "create", "namespace", k.namespace,
	)
	return common.ExecCmd(cmd)
}

/*
Delete a namespace in k8s
*/
func (k *Kubernetes) DeleteNamespace() error {
	defer common.TimeCurrentFn()()

	cmd := exec.Command(
		"kubectl", "delete", "namespace", k.namespace,
	)
	return common.ExecCmd(cmd)
}

/*
Apply all mstk kubernetes files to the namespace
*/
func (k *Kubernetes) ApplyMSTK() error {
	dir, err := common.GetMSTKDir()
	if err != nil {
		return err
	}

	return k.Apply(dir + "k8s/")
}

/*
Delete all mstk kubernetes files from the namespace
*/
func (k *Kubernetes) DeleteMSTK() error {
	dir, err := common.GetMSTKDir()
	if err != nil {
		return err
	}

	return k.Delete(dir + "k8s/")
}

// PRIVATE FUNCTIONS
// ========================================================================
