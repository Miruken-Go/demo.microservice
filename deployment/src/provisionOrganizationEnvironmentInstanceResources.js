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

        logging.header("Deploying Organization Instance Resources")

        //Clean up from any deleted resources
        await az.deleteOrphanedApplicationSecurityPrincipals()

        //Resources Groups
        await az.createResourceGroup(organization.resourceGroups.instance, organization.location)

        const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
        const bicepFile                   = path.join(__dirname, 'bicep/organizationEnvironmentInstanceResources.bicep')

        const applications = organization.applications.map(a => {
            return { 
                name:             a.name, 
                containerAppName: a.containerAppName, 
                secrets:          a.secrets
            }
        })

        const params = JSON.stringify({ 
            prefix:                      { value: organization.resourceGroups.instance },
            containerRepositoryName:     { value: organization.containerRepositoryName },
            location:                    { value: organization.location },
            keyVaultResourceGroup:       { value: organization.resourceGroups.common },
            keyVaultName:                { value: organization.keyVaultName },
            containerRepositoryPassword: { value: containerRepositoryPassword },
            applications:                { value: applications },
        })

        const results = await bash.json(`
            az deployment group create                                   \
                --template-file  ${bicepFile}                            \
                --subscription   ${variables.subscriptionId}             \
                --resource-group ${organization.resourceGroups.instance} \
                --mode           complete                                \
                --parameters     '${params}'                             \
        `)

        logging.printObject("Bicep Outputs", results.properties.outputs)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
