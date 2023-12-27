import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as logging     from '#infrastructure/logging.js'
import * as gh          from '#infrastructure/gh.js'
import { handle }       from '#infrastructure/handler.js'
import { variables }    from '#infrastructure/envVariables.js'
import { organization } from './config.js'

variables.requireEnvVariables([
    'subscriptionId',
])

handle(async () => {
    logging.printEnvironmentVariables(variables)
    logging.printOrganization(organization)

    logging.header("Deploying Organization Instance Resources")

    //Clean up from any deleted resources
    await az.deleteOrphanedApplicationSecurityPrincipals()

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.instance, organization.location)

    const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
    const bicepFile                   = new URL('bicep/organizationEnvironmentInstanceResources.bicep', import.meta.url).pathname

    const applications = organization.applications.map(a => {
        return { 
            name:             a.name, 
            containerAppName: a.containerAppName, 
            secrets:          a.secrets
        }
    })

    const params = JSON.stringify({ 
        prefix:                      { value: organization.resourceGroups.instance },
        containerRepositoryName:     { value: organization.containerRepositoryName },
        location:                    { value: organization.location },
        keyVaultResourceGroup:       { value: organization.resourceGroups.common },
        keyVaultName:                { value: organization.keyVaultName },
        containerRepositoryPassword: { value: containerRepositoryPassword },
        applications:                { value: applications },
    })

    const results = await bash.json(`
        az deployment group create                                     \
            --name           OrgInstance${Math.floor(Date.now()/1000)} \                   
            --template-file  ${bicepFile}                              \
            --subscription   ${variables.subscriptionId}               \
            --resource-group ${organization.resourceGroups.instance}   \
            --mode           complete                                  \
            --parameters     '${params}'                               \
    `)

    logging.printObject("Bicep Outputs", results.properties.outputs)

    await gh.sendRepositoryDispatch(`provisioned-organization-environment-instance-resources`, {
        env:      organization.env,
        instance: organization.instance,
    })
})
