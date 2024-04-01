import {
    handle,
    EnvVariables,
    EnvSecrets,
    logging,
    AZ,
    bash
} from 'ci.cd'
import { organization } from './config.js'


handle(async () => {
    const variables = new EnvVariables()
        .required([
            'tenantId',
            'subscriptionId',
            'deploymentPipelineClientId',
            'deploymentPipelineClientSecret',
            'location',
        ])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
        ])
        .secrets
    logging.printSecrets(secrets)

    logging.printDomain(organization)

    logging.header("Deploying Organization Global Resources")

    const az = new AZ({
        tenantId:                       variables.tenantId,
        subscriptionId:                 variables.subscriptionId,
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
    })

    //Provider Registrations
    await az.registerAzureProvider('Microsoft.AzureActiveDirectory')
    await az.registerAzureProvider('Microsoft.App')
    await az.registerAzureProvider('Microsoft.OperationalInsights')

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.global, organization.location, {})

    const bicepFile = new URL('bicep/organizationGlobalResources.bicep', import.meta.url).pathname

    await bash.json(`
        az deployment group create                                                     \
            --name           organizationGlobalResources${Math.floor(Date.now()/1000)} \
            --template-file  ${bicepFile}                                              \
            --subscription   ${variables.subscriptionId}                               \
            --resource-group ${organization.resourceGroups.global}                     \
            --mode complete                                                            \
            --parameters                                                               \
                containerRepositoryName=${organization.containerRepository.name}       \
                location=${organization.location}                                      \
    `)
})
