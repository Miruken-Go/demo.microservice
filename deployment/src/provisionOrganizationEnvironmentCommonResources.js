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
    logging.printDomain(organization)

    logging.header("Deploying Organization Common Resources")

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.common, organization.location)
    await az.createResourceGroup(organization.resourceGroups.manual, organization.location)

    const bicepFile = new URL('bicep/organizationEnvironmentCommonResources.bicep', import.meta.url).pathname

    await bash.json(`
        az deployment group create                                                                \
            --name           organizationEnvironmentCommonResources${Math.floor(Date.now()/1000)} \
            --template-file  ${bicepFile}                                                         \
            --subscription   ${variables.subscriptionId}                                          \
            --resource-group ${organization.resourceGroups.common}                                \
            --mode complete                                                                       \
            --parameters                                                                          \
                prefix=${organization.resourceGroups.common}                                      \
                keyVaultName=${organization.keyVault.name}                                         \
                location=${organization.location}                                                 \
    `)

    await gh.sendRepositoryDispatch(`provisioned-organization-environment-common-resources`, {
        env: organization.env
    })
})
