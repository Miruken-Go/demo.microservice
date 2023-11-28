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

        logging.header("Deploying Organization Common Resources")

        //Resources Groups
        await az.createResourceGroup(organization.resourceGroups.common, organization.location)
        await az.createResourceGroup(organization.resourceGroups.manual, organization.location)

        const bicepFile = path.join(__dirname, 'bicep/organizationEnvironmentCommonResources.bicep')

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

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
