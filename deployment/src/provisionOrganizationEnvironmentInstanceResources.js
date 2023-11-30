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

        logging.header("Deploying Organization Instance Resources")

        //Clean up from any deleted resources
        await az.deleteOrphanedApplicationSecurityPrincipals()

        //Resources Groups
        await az.createResourceGroup(organization.resourceGroups.instance, organization.location)

        const containerRepositoryPassword = await az.getAzureContainerRepositoryPassword(organization.containerRepositoryName)
        const bicepFile                   = new URL('bicep/organizationEnvironmentInstanceResources.bicep', import.meta.url).pathname

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
