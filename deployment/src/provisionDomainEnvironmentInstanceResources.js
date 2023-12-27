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

    const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
    const bicepFile                   = new URL('bicep/domainEnvironmentInstanceResources.bicep', import.meta.url).pathname
    const tags                        = `organization=${organization.name}`

    for (const domain of organization.domains) {

        //Resources Groups
        await az.createResourceGroup(domain.resourceGroups.instance, domain.location, tags)

        const applications = domain.applications.map(a => {
            return { 
                name:             a.name, 
                containerAppName: a.containerAppName, 
                secrets:          a.secrets
            }
        })

        const params = JSON.stringify({ 
            containerRepositoryPassword: { value: containerRepositoryPassword },
            containerRepositoryName:     { value: organization.containerRepositoryName },
            keyVaultResourceGroup:       { value: organization.resourceGroups.common },
            keyVaultName:                { value: organization.keyVaultName },
            prefix:                      { value: domain.resourceGroups.instance },
            location:                    { value: domain.location },
            applications:                { value: applications },
        })

        const results = await bash.json(`
            az deployment group create                                        \
                --name           DomainInstance${Math.floor(Date.now()/1000)} \
                --template-file  ${bicepFile}                                 \
                --subscription   ${variables.subscriptionId}                  \
                --resource-group ${domain.resourceGroups.instance}            \
                --mode           complete                                     \
                --parameters     '${params}'                                  \
        `)

        logging.printObject("Bicep Outputs", results.properties.outputs)
    }

    await gh.sendRepositoryDispatch(`provisioned-domain-environment-instance-resources`, {
        env:      organization.env,
        instance: organization.instance,
    })
})
