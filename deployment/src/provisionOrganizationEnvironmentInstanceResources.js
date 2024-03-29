import * as az                          from '#infrastructure/az.js'
import * as bash                        from '#infrastructure/bash.js'
import * as logging                     from '#infrastructure/logging.js'
import * as gh                          from '#infrastructure/gh.js'
import { handle }                       from '#infrastructure/handler.js'
import { variables }                    from '#infrastructure/envVariables.js'
import { getImageTagForActiveRevision } from '#infrastructure/containerApp.js'
import { organization }                 from './config.js'

variables.requireEnvVariables([
    'subscriptionId',
])

handle(async () => {
    logging.printEnvironmentVariables(variables)
    logging.printDomain(organization)

    logging.header("Deploying Organization Instance Resources")

    //Clean up from any deleted resources
    await az.deleteOrphanedApplicationSecurityPrincipals()

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.instance, organization.location)

    const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepository.name)
    const bicepFile                   = new URL('bicep/organizationEnvironmentInstanceResources.bicep', import.meta.url).pathname

    const applications = []
    for(const a of organization.applications) {
        const imageTag = (await getImageTagForActiveRevision(a)) || 'default'
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
        keyVaultName:                { value: organization.keyVault.name },
        containerRepositoryPassword: { value: containerRepositoryPassword },
        applications:                { value: applications },
    })

    const results = await bash.json(`
        az deployment group create                                                                  \
            --name           organizationEnvironmentInstanceResources${Math.floor(Date.now()/1000)} \
            --template-file  ${bicepFile}                                                           \
            --subscription   ${variables.subscriptionId}                                            \
            --resource-group ${organization.resourceGroups.instance}                                \
            --mode           complete                                                               \
            --parameters     '${params}'                                                            \
    `)

    logging.printObject("Bicep Outputs", results.properties.outputs)

    await gh.sendRepositoryDispatch(`provisioned-organization-environment-instance-resources`, {
        env:      organization.env,
        instance: organization.instance,
    })
})
