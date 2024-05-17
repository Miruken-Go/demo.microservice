# build-image

Docker container for running all cicd scripts.
The container is published to github packages.

## Run the "build-and-push-deployment-image" pipeline to build a new image

  Search and replace the container version `ghcr.io/miruken-go/demo.microservice.build:1713989236` in the `demo.microservice` folder.

## To Develop Locally

Build the Docker Container

    docker build -t demo.microservice.build:local .

Run the Docker Container interactively

    docker run -it --rm -v $(pwd):/build demo.microservice.build:local

Run the container with required environment variables.
I do this from a bash script at the same level as the demo.microservice folder
so that there is no risk of secrets being accidentally checked in.

    docker run -it --rm                                                    \
        -v $(pwd)/demo.microservice:/build                                 \
        -v /var/run/docker.sock:/var/run/docker.sock                       \
        -w /build                                                          \
        -e tenantId=<tenantId>                                             \
        -e subscriptionId=<subscriptionId>                                 \
        -e deploymentPipelineClientId=<deploymentPipelineClientId>         \
        -e deploymentPipelineClientSecret=<deploymentPipelineClientSecret> \
        -e GH_TOKEN=<GH_TOKEN>                                             \
        -e repository='Miruken-Go/demo.microservice'                       \
        -e repositoryOwner='Miruken-Go'                                    \
        -e env=dev                                                         \
        -e instance=<optional>                                             \
        demo.microservice.deployment:local                                 \
        /bin/bash                                                          \


