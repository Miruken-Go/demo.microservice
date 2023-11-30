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

        logging.header("Deploying Organization Global Resources")

        //Provider Registrations
        await az.registerAzureProvider('Microsoft.AzureActiveDirectory')
        await az.registerAzureProvider('Microsoft.App')
        await az.registerAzureProvider('Microsoft.OperationalInsights')

        //Resources Groups
        await az.createResourceGroup(organization.resourceGroups.global, organization.location)

        const bicepFile = new URL('bicep/organizationGlobalResources.bicep', import.meta.url).pathname

        await bash.json(`
            az deployment group create                                              \
                --template-file  ${bicepFile}                                       \
                --subscription   ${variables.subscriptionId}                        \
                --resource-group ${organization.resourceGroups.global}              \
                --mode complete                                                     \
                --parameters                                                        \
                    containerRepositoryName=${organization.containerRepositoryName} \
                    location=${organization.location}                               \
        `)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
