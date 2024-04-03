# build-image

Docker container for running all cicd scripts.

### Run the "build-and-push-deployment-image" pipeline

Get the deployment tag from the container and update the 
other pipeline files with the latest tag.

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


