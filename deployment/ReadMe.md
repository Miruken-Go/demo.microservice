# demo.microservice.deployment

## Initial Workflow when starting the project

### Manually create deployment permissions in Azure

Create a service principal
https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal

    Azure AD > App Registrations > Add
    Name it DeploymentPipeline
    Certificates&Secrets > New Client Secret

Give DeploymentPipeline Graph Api permissions
    Azure AD > App Registrations > API permissions
        Add a permission > Microsoft Graph 
            Directory.Raed.All

Give DeploymentPipeline permissions on the subscription where the app will create resources

    Subscriptions > Access control (IAM)
    Add > Add role assignement > Owner
    Add > Add role assignment  > Key Vault Secrets User

### Add secrets and variables in github

Settings > Security > Secrets and Variables > Actions

Secrets

* DEPLOYMENT_PIPELINE_CLIENT_SECRET
* WORKFLOW_GH_TOKEN

`WORKFLOW_GH_TOKEN` should be a github personal access token that has permissions to edit workflows
    
Variables

* TENANT_ID
* SUBSCRIPTION_ID
* DEPLOYMENT_PIPELINE_CLIENT_ID


### Run the "build-and-push-deployment-image" pipeline

Get the deployment tag from the container and update the 
other pipeline files with the latest tag.

### Run the "deploy-global-resources" pipeline

### Push the initial image to the shared Azure Container Repository

    TAG=initial
    IMAGE_NAME="teamsrvglobal.azurecr.io/teamsrv:$TAG"; echo $IMAGE_NAME
    docker build --build-arg application_version=$TAG -t $IMAGE_NAME demo.microservice/teamsrv 
    az acr login -n teamsrvglobal
    docker push $IMAGE_NAME

### Run the "deploy-environment" pipeline

## To Develop Locally

Build the Docker Container

    docker build -t demo.microservice.build:local .

Run the Docker Container interactively

    docker run -it --rm -v $(pwd):/build demo.microservice.build:local

Execute the build

    docker run -it --rm                                                    \
        -v $(pwd):/build                                                   \
        -e tenantId=<tenantId>                                             \
        -e subscriptionId=<subscriptionId>                                 \
        -e deploymentPipelineClientId=<deploymentPipelineClientId>         \
        -e deploymentPipelineClientSecret=<deploymentPipelineClientSecret> \
        -e env=dev                                                         \
        -e instance=craig                                                  \
        demo.microservice.deployment:local                                 \
        node /build/src


