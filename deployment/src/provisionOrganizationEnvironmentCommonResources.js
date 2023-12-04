import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as logging     from '#infrastructure/logging.js'
import { handle }       from '#infrastructure/handler.js'
import { variables }    from '#infrastructure/envVariables.js'
import { organization } from './config.js'

variables.requireEnvVariables([
    'subscriptionId',
])

handle(async () => {
    logging.printEnvironmentVariables(variables)
    logging.printOrganization(organization)

    logging.header("Deploying Organization Common Resources")

    //Resources Groups
    await az.createResourceGroup(organization.resourceGroups.common, organization.location)
    await az.createResourceGroup(organization.resourceGroups.manual, organization.location)

    const bicepFile = new URL('bicep/organizationEnvironmentCommonResources.bicep', import.meta.url).pathname

    await bash.json(`
        az deployment group create                                 \
            --template-file  ${bicepFile}                          \
            --subscription   ${variables.subscriptionId}           \
            --resource-group ${organization.resourceGroups.common} \
            --mode complete                                        \
            --parameters                                           \
                keyVaultName=${organization.keyVaultName}          \
                location=${organization.location}                  \
    `)

    await gh.sendRepositoryDispatch(`provisioned-organization-environment-common-resources`, {
        env: organization.env
    })
})
