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

### Deploy your microservices
Now your system is ready to use mstk. Inside your project run

```
mstk deploy
```

This command will read `Service.toml` in the current directory and according to it will build and push your load balancers and microservices

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
