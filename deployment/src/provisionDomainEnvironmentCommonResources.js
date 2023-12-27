import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as gh          from '#infrastructure/gh.js'
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

    logging.header("Deploying Domain Common Resources")

    const tags = `organization=${organization.name}`
    for (const domain of organization.domains) {
        //Resources Groups
        await az.createResourceGroup(domain.resourceGroups.common, domain.location, tags)

        const bicepFile = new URL('bicep/domainEnvironmentCommonResources.bicep', import.meta.url).pathname

        const results = await bash.json(`
            az deployment group create                                                          \
                --name           domainEnvironmentCommonResources${Math.floor(Date.now()/1000)} \
                --template-file  ${bicepFile}                                                   \
                --subscription   ${variables.subscriptionId}                                    \
                --resource-group ${domain.resourceGroups.common}                                \
                --mode complete                                                                 \
                --parameters                                                                    \
                    location=${domain.location}                                                 \
        `)
            
        logging.printObject("Bicep Outputs", results.properties.outputs)
    }

    await gh.sendRepositoryDispatch(`provisioned-domain-environment-common-resources`, {
        env: organization.env
    })
})
