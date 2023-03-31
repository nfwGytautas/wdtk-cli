# MSTK
Microservice Toolkit is a suite of tools whose main purpose is to make the creation of distributed microservice systems easier.

## How to use

Run these commands to setup MSTK on your local machine (kubectl via minikube needs to be set up already and go needs to be installed)

### Setup
```
git clone github.com/nfwGytautas/mstk
cd cli
go build -o mstk *.go
mstk setup
```

### User startup

MSTK CLI has some utility commands prepared for you to make development easier.

First you need a MSTK project which you can scaffold using

```
mstk template <package> <project>
```

The ```<package>``` keyword is the url to your project (the string that will be written in go.mod files)

This will create a MSTK template project with the name
```<project>``` whose tree structure looks like this:

```
.
├── balancers
├── go.work
├── mstk_project.toml
└── services
```

Now to create a service you can run the CLI command

```
mstk service <name>
```

This command wil automatically create a service and a load balancer directory for you and modify the necessary files to make it work seamlessly. Altho this can be done by hand using the mstk CLI is more convenient. After running the command you should have something like

```
.
├── go.work
├── mstk_project.toml
└── services
    └── name_of_service
        ├── Service.toml
        ├── balancer
        │   ├── go.mod
        │   └── main.go
        └── service
            ├── go.mod
            └── main.go
```


###  Deploy your microservices
Once you have developed your microservice you can deploy it with

```
mstk deploy
```

This command will read `mstk_project.toml` in the current directory and according to it will build and push your load balancers and microservices. The command is smart and will only deploy if it detects changes in microservices.

Alternatively you can use

```
mstk deploy <service>
```

To deploy a specific service

### Shutdown
Once you are done you can shutdown mstk using

```
mstk clean
```

This command will cleanup all mstk related kubernetes pods aswell as any other artifacts of the setup command


## Project structure

```
.
├── README.md
└── gomods      - Directory containing the modules in go
```
