import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as logging     from '#infrastructure/logging.js'
import { variables }    from '#infrastructure/envVariables.js'
import { organization } from './config.js'

variables.requireEnvVariables([
    'subscriptionId',
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        logging.header("Deploying Domain Common Resources")

        const tags = `organization=${organization.name}`
        for (const domain of organization.domains) {
            //Resources Groups
            await az.createResourceGroup(domain.resourceGroups.common, domain.location, tags)

            const bicepFile = new URL('bicep/domainEnvironmentCommonResources.bicep', import.meta.url).pathname

            const results = await bash.json(`
                az deployment group create                                 \
                    --template-file  ${bicepFile}                          \
                    --subscription   ${variables.subscriptionId}           \
                    --resource-group ${domain.resourceGroups.common}       \
                    --mode complete                                        \
                    --parameters                                           \
                        location=${domain.location}                  \
            `)
                
            logging.printObject("Bicep Outputs", results.properties.outputs)
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
