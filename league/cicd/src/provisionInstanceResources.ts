import {
    handle,
    EnvVariables,
    EnvSecrets,
    logging,
    AZ,
    bash,
    containerApp,
    Domain
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
            'env'
        ])
        .optional(['instance'])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
        ])
        .secrets
    logging.printSecrets(secrets)

    logging.header("Deploying League Instance Resources")

    const az = new AZ({
        tenantId:                       variables.tenantId,
        subscriptionId:                 variables.subscriptionId,
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
    })

    //Clean up from any deleted resources
    await az.deleteOrphanedApplicationSecurityPrincipals()

    const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.resources.containerRepository.name)
    const bicepFile                   = join(__dirname, 'bicep/instanceResources.bicep')

    //Resources Groups
    const domain: Domain = organization.getDomainByName('league')

    const tags = {
        organization: organization.name,
        domain:       domain.name,
        env:          variables.env,
        instance:     variables.instance ?? ''
    }
 
    await az.createResourceGroup(domain.resourceGroups.instance, domain.location, tags)

    const applications: applicationsParam[] = []
    for(const a of domain.applications) {
        const imageTag = (await containerApp.getImageTagForActiveRevision(a)) || 'default'
        applications.push({ 
            name:             a.name, 
            containerAppName: a.containerAppName, 
            secrets:          a.secrets,
            imageTag,
        })
    }

    const params = JSON.stringify({ 
        containerRepositoryPassword: { value: containerRepositoryPassword },
        prefix:                      { value: domain.resourceGroups.instance },
        location:                    { value: domain.location },
        containerRepositoryName:     { value: organization.resources.containerRepository.name },
        applications:                { value: applications },
        tags:                        { value: tags },
    })

    const results = await bash.json(`
        az deployment group create                                           \
            --name           instanceResources${Math.floor(Date.now()/1000)} \
            --template-file  ${bicepFile}                                    \
            --subscription   ${variables.subscriptionId}                     \
            --resource-group ${domain.resourceGroups.instance}               \
            --mode           complete                                        \
            --parameters     '${params}'                                     \
    `)

    logging.printObject("Bicep Outputs", results.properties.outputs)
})

interface applicationsParam {
    name:             string 
    containerAppName: string 
    imageTag:         string
    secrets:          string[]
}
