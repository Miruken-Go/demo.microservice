const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const { variables }    = require('./infrastructure/envVariables')
const { organization } = require('./config');
const path             = require('path')

variables.require([
    'subscriptionId',
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        logging.header("Deploying Organization Instance Resources")

        //Resources Groups
        await az.createResourceGroup(organization.resourceGroups.instance, organization.location)

        const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
        const bicepFile                   = path.join(__dirname, 'bicep/organizationInstanceResources.bicep')

        const params = JSON.stringify({ 
            prefix:                      organization.resourceGroups.instance,
            containerRepositoryName:     organization.containerRepositoryName,
            location:                    organization.location,
            keyVaultResourceGroup:       organization.resourceGroups.common,
            keyVaultName:                organization.keyVaultName,
            containerRepositoryPassword: containerRepositoryPassword,
            // applications: organization.applications.map(x => {
            //     return { 
            //         name:    x.name, 
            //         secrets: x.secrets
            //     }
            // })
        })

        const results = await bash.json(`
            az deployment group create                                              \
                --template-file  ${bicepFile}                                       \
                --subscription   ${variables.subscriptionId}                        \
                --resource-group ${organization.resourceGroups.instance}            \
                --mode complete                                                     \
                --parameters                                                        \
                    prefix=${organization.resourceGroups.instance}                  \
                    containerRepositoryName=${organization.containerRepositoryName} \
                    location=${organization.location}                               \
                    keyVaultResourceGroup=${organization.resourceGroups.common}     \
                    keyVaultName=${organization.keyVaultName}                       \
                    containerRepositoryPassword=${containerRepositoryPassword}      \
        `)

        logging.header("Container App Urls")
        console.log(results.properties.outputs.containerAppUrls.value)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }



    
}

main()
