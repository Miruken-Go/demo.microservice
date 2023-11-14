const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const { variables }    = require('./infrastructure/envVariables')
const { organization } = require('./config');
const path             = require('path')

variables.require([
    'tenantId',
    'subscriptionId',
    'deploymentPipelineClientId',
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        logging.header("Deploying Organization Common Resources")

        //Resources Groups
        //await az.createResourceGroup(organization.resourceGroups.common, organization.location)
        //await az.createResourceGroup(organization.resourceGroups.manual, organization.location)
        
        //await az.createResourceGroup(organization.resourceGroups.stable, organization.location)
        if ( variables.instance) {
            await az.createResourceGroup(organization.resourceGroups.instance, organization.location)
        }

        logging.header("Deploying OrganizationGlobalResources Arm Template")
        const bicepFile = path.join(__dirname, 'bicep/organizationGlobalResources.bicep')

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
