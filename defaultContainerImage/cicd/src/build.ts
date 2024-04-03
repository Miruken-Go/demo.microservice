import { 
    handle,
    EnvVariables,
    EnvSecrets,
    bash,
    logging,
    AZ,
    Git,
    Application
} from 'ci.cd'

import {organization} from './domains'

handle(async () => {
   const variables = new EnvVariables()
        .required([
            'deploymentPipelineClientId',
            'subscriptionId',
            'tenantId'
        ])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
            'GH_TOKEN'
        ])
        .secrets
    logging.printSecrets(secrets)

    logging.header("Building default-container-image for all applications")

    const version   = `v${Math.floor(Date.now()/1000)}`.trim()
    const tag       = `default-container-image/${version}`
    const imageName = 'default-container-image:latest' 

    console.log(`version: [${version}]`)
    console.log(`tag:     [${tag}]`)

    await bash.execute(`
        cd ../app
        docker build                \
            -t ${imageName}         \
            .                       \
    `)

    await new AZ({
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret,
        subscriptionId:                 variables.subscriptionId,
        tenantId:                       variables.tenantId
    }).loginToACR(organization.resources.containerRepository.name)

    //Push the default container for all the configured apps
    for (const app of organization.applications) {
        await tagContainerImageAndPush(imageName, app)
    }
    for(const domain of organization.domains) {
        for(const app of domain.applications) {
            await tagContainerImageAndPush(imageName, app)
        }
    }

    await new Git(secrets.GH_TOKEN)
        .tagAndPush(tag)
})

async function tagContainerImageAndPush(imageName:string, app: Application) {
    const appImage = `${app.imageName}:default`
    console.log(`imageName: [${appImage}]`)
    await bash.execute(`
        docker tag ${imageName} ${appImage}
        docker push ${appImage}
    `)
}