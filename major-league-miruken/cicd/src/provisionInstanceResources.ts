import {
    handle,
    EnvVariables,
    EnvSecrets,
    logging,
    AZ,
    bash,
    containerApp
} from 'ci.cd'

import { organization } from './domains'
import { join } from 'node:path'


handle(async () => {
    const variables = new EnvVariables()
        .required([
            'tenantId',
            'subscriptionId',
            'deploymentPipelineClientId',
            'deploymentPipelineClientSecret',
        ])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
        ])
        .secrets
    logging.printSecrets(secrets)

    logging.header("Deploying Organization Instance Resources")

    const az = new AZ({
        tenantId:                       variables.tenantId,
        subscriptionId:                 variables.subscriptionId,
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
    })

    //Clean up from any deleted resources
    await az.deleteOrphanedApplicationSecurityPrincipals()

    const tags = {
        organization: organization.name,
        domain:       organization.name,
        env:          variables.env,
        instance:     variables.instance ?? ''
    }

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.instance, organization.location, tags)

    const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepository.name)
    const bicepFile = join(__dirname, 'bicep/instance.bicep')

    const applications: applicationsParam[] = []
    for(const a of organization.applications) {
        const imageTag = (await containerApp.getImageTagForActiveRevision(a)) || 'default'
        applications.push({ 
            name:             a.name, 
            containerAppName: a.containerAppName, 
            secrets:          a.secrets,
            imageTag,
        })
    }

    const params = JSON.stringify({ 
        prefix:                      { value: organization.resourceGroups.instance },
        containerRepositoryName:     { value: organization.containerRepository.name },
        location:                    { value: organization.location },
        keyVaultResourceGroup:       { value: organization.resourceGroups.common },
        keyVaultName:                { value: organization.resources.keyVault.name },
        containerRepositoryPassword: { value: containerRepositoryPassword },
        applications:                { value: applications },
        tags:                        { value: tags },
    })

    const results = await bash.json(`
        az deployment group create                                                                  \
            --name           instanceResources${Math.floor(Date.now()/1000)} \
            --template-file  ${bicepFile}                                                           \
            --subscription   ${variables.subscriptionId}                                            \
            --resource-group ${organization.resourceGroups.instance}                                \
            --mode           complete                                                               \
            --parameters     '${params}'                                                            \
    `)

    logging.printObject("Bicep Outputs", results.properties.outputs)
})

interface applicationsParam {
    name:             string 
    containerAppName: string 
    imageTag:         string
    secrets:          string[]
}