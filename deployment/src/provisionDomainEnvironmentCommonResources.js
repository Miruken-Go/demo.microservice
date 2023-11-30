const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const { variables }    = require('./infrastructure/envVariables')
const { organization } = require('./config');
const path             = require('path')

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
