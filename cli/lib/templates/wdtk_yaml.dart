part of wdtk_templates;

/// Returns a wdtk.yaml string template
String templateWdtkYaml(String domain, String projectName) {
  return """
# Generic project settings
package: $domain
name: $projectName

# List of variables, these can be reference anywhere using \${variableName}
# Variables can have arguments defined as \${#nameOfArgument}
# When referencing them one can use \${variableName, nameOfArgument: something}, to insert into the placeholders
aliases:
    connectionOptions: "?charset=utf8mb4&parseTime=True&loc=Local"
    databaseStringBase: "user:password@tcp(\${#database_ip}:3306)/\${#database}\${connectionOptions}"
    apiKey: API_KEY_GOES_HERE

# Configure deployments for services in the deployments section
deployments:
  - name: dev
    # The 'defaults' entry contains the default deployment settings for each service
    # 'defaults' has to have all values defined unless they are explicitly set by hand for each service
    # In this example since port is set in all services it doesn't need to be specified in the defaults section
    defaults:
      ip: 127.0.0.1
      # Referencing an alias defined before
      apiKey: \${apiKey}
      # Here we use some standard aliases starting and ending with '__' these are set by wdtk
      # and can be used in any context as aliases can
      dir: \${__HOME__}/\${__PACKAGE__}/\${__DEPLOYMENT__}/\${__SERVICE__}
    # Deployments can also describe aliases these can be accessed with \${::variable} notation
    aliases:
      # By using the '...' keyword we forward the remaining arguments, e.g. database to whom ever is going to use the alias
      databaseString: "\${databaseStringBase, database_ip: 127.0.0.1, ...}"
    # The 'settings' entry can contain service names, which then can also contains all the configuration data
    settings:
      Gateway:
        port: 8090
      Authentication:
        port: 8091
      Http-Server:
        port: 8080
      ExampleService:
        port: 8100

# Describe the frontends
frontend:
  platforms:
    # Write all the used frontend platforms
    - type: web
      toolchain: flutter

# Services array is used to define the services inside this package
services:
  - name: Gateway
    source:
      # We are using a remote go service that is accessible through git
      type: git
      remote: github.com/nfwGytautas/wdtk-services/gateway
      language: go
    options:
      gateway: true

  - name: Authentication
    source:
      type: git
      remote: github.com/nfwGytautas/wdtk-services/authentication
      language: go
    # Config key can be used for additional configuration options these will be stored inside the generated service config files
    config:
      # Using a parameterized alias for ease of access, here we are accessing a deployment alias because of '::' prefix
      connectionString: "\${::databaseString, database: auth}"

  # Simple http-server for the web frontend
  - name: Http-Server
    source:
      type: git
      remote: github.com/nfwGytautas/wdtk-services/http-server
      language: go
    config:
      # __WEB_DEPLOYMENT_DIR__ is a unique alias for accessing the deployment dir of the web frontend
      htmlDirectory: "\${__DEPLOYMENT_DIR__, service: Web}"

  # Describe services here
  - name: ExampleService
    source:
      type: src
      language: go
""";
}