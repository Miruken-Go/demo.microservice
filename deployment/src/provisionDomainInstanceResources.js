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

        const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
        const bicepFile                   = path.join(__dirname, 'bicep/domainInstanceResources.bicep')
        const tags                        = `organization=${organization.name}`

        for (const domain of organization.domains) {

            //Resources Groups
            await az.createResourceGroup(domain.resourceGroups.instance, domain.location, tags)

            const applications = domain.applications.map(a => {
                return { 
                    name:             a.name, 
                    containerAppName: a.containerAppName, 
                    secrets:          a.secrets
                }
            })

            const params = JSON.stringify({ 
                containerRepositoryPassword: { value: containerRepositoryPassword },
                containerRepositoryName:     { value: organization.containerRepositoryName },
                keyVaultResourceGroup:       { value: organization.resourceGroups.common },
                keyVaultName:                { value: organization.keyVaultName },
                prefix:                      { value: domain.resourceGroups.instance },
                location:                    { value: domain.location },
                applications:                { value: applications },
            })

            const results = await bash.json(`
                az deployment group create                                   \
                    --template-file  ${bicepFile}                            \
                    --subscription   ${variables.subscriptionId}             \
                    --resource-group ${domain.resourceGroups.instance} \
                    --mode           complete                                \
                    --parameters     '${params}'                             \
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
