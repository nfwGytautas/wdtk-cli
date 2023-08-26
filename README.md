# WDTK
WebDev-ToolKit. This is a collection of libraries, tools for creating a web related technologies. The idea is to provide a simple way to create easily extendable microservice backends and multi target frontends (app, webpage, desktop app).

## Usage
Most of the functionality happens through the cli and a yaml based config file.

### Setup
The recommended way to setup is by cloning the repository

On UNIX:
```
git clone https://github.com/nfwGytautas/webdev-tk.git
cd ./cli/
go mod tidy
go build -o wdtk ./
mv wdtk /bin/
```

The last move command is not mandatory but its easier to use then

### Creating a project
You can create a new project with
```
cd PROJECT_DIRECTORY
wdtk init -n NAME_OF_PROJECT
```

After running this command you will get a very basic file and directory structure for the wdtk project and it will automatically create a `wdtk.yml `file. For more information about the configuration file you can find [HERE](documentation/CONFIGURATION_FILE.md)

### Scaffold command
When ever you change, clone or pull changes of an existing wdtk project (e.g. the configuration file) you will want to run `wdtk scaffold` this command will check the configuration file in will make sure that everything is correctly configured and working. It will report any warnings, errors, other useful information about the project in the command line. More information [HERE](documentation/SCAFFOLD_COMMAND.md)

### Deploy command
When it is finally time to deploy the cli provided a command `wdtk deploy` this command expects a deployment target as the first parameter and the following arguments shall be either all or a list of services to deploy. The command will then execute deployment for the specified configuration and will write a log file inside `.wdtk/logs/` directory in your project root. More information [HERE](documentation/DEPLOY_COMMAND.md)
