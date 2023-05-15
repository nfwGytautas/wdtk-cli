package templates

// PUBLIC TYPES
// ========================================================================

/*
Data of service template
*/
type UNIXDeployData struct {
	RootDir      string
	ServiceName  string
	BalancerLang string
	ServiceLang  string
}

// PUBLIC FUNCTIONS
// ========================================================================

/*
Template for deploy scripts for linux
*/
const UnixDeployTemplate = `
#/bin/bash

SERVICE_NAME="{{.ServiceName}}"

BALANCER_LANGUAGE="{{.BalancerLang}}"
BALANCER_SRC_DIR="{{.RootDir}}/services/{{.ServiceName}}/balancer/"

SERVICE_LANGUAGE="{{.ServiceLang}}"
SERVICE_SRC_DIR="{{.RootDir}}/services/{{.ServiceName}}/service/"

OUT_DIR="{{.RootDir}}/deploy/bin/"

echo "Start of '$SERVICE_NAME'"
echo

# build balancer
if [ $BALANCER_LANGUAGE = "go" ]
then
    echo Building balancer

    OBJ_FILE_NAME=$SERVICE_NAME'_balancer'

    cd $BALANCER_SRC_DIR

    echo Running go tidy
    go mod tidy

    echo Downloading dependencies
    go get ./

    echo Building
    go build -o $OUT_DIR$OBJ_FILE_NAME .

    echo
fi

# build service
if [ $SERVICE_LANGUAGE = "go" ]
then
    echo Building service

    OBJ_FILE_NAME=$SERVICE_NAME'_service'

    cd $SERVICE_SRC_DIR

    echo Running go tidy
    go mod tidy

    echo Downloading dependencies
    go get ./

    echo Building
    go build -o $OUT_DIR$OBJ_FILE_NAME .
fi

echo "End of '$SERVICE_NAME'"

`
