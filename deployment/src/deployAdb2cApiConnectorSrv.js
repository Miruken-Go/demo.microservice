const logging            = require('./infrastructure/logging');
const az                 = require('./infrastructure/az');
const bash               = require('./infrastructure/bash')
const { variables }      = require('./infrastructure/envVariables')
const { organization }   = require('./config');

variables.requireEnvVariables([
    'tag'
])

variables.requireEnvFileVariables([
    'authorizationServiceUsername',
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        const application = organization.getApplicationByName("adb2c-api-connector-srv")

        logging.header(`Deploying ${application.name}`)

        const envVars = [
            `Login__Adb2c__0__Module='login.pwd'`,
            `Login__Adb2c__0__Options__Credentials__0__Username='${variables.authorizationServiceUsername}'`,
            `Login__Adb2c__0__Options__Credentials__0__Password='secretref:authorization-service-password'`,
        ]

        await az.login()

        //https://learn.microsoft.com/en-us/cli/azure/containerapp?view=azure-cli-latest#az-containerapp-update
        //Create the new revision
        const now            = `${Math.floor(Date.now()/1000)}`.trim()
        const revisionSuffix = `${variables.tag}-${now}`
        await bash.execute(`
            az containerapp update                                \
                -n ${application.containerAppName}                \
                -g ${application.resourceGroups.instance}         \
                --image ${application.imageName}:${variables.tag} \
                --container-name ${application.name}              \
                --revision-suffix ${revisionSuffix}               \
                --replace-env-vars ${envVars.join(' ')}           \
        `)

        const revisions = await bash.json(`
            az containerapp revision list                 \
                -n ${application.containerAppName}        \
                -g ${application.resourceGroups.instance} \
        `)

        const revisionsToActivate   = revisions.filter(r => r.name.includes(revisionSuffix))
        const revisionsToDeactivate = revisions.filter(r => !r.name.includes(revisionSuffix))

        //You must have an active revision before deactivating the rest
        for (const revision of revisionsToActivate) {
            if (revision.properties.active !== true) {
                console.log(`Activating: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision activate             \
                        -n ${application.containerAppName}        \
                        -g ${application.resourceGroups.instance} \
                        --revision ${revision.name}               \
                `)
            }

            await bash.execute(`
                az containerapp ingress traffic set           \
                    -n ${application.containerAppName}        \
                    -g ${application.resourceGroups.instance} \
                    --revision-weight ${revision.name}=100    \
            `)
        }

        for (const revision of revisionsToDeactivate) {
            if (revision.properties.active === true) {
                console.log(`Dectivate: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision deactivate           \
                        -n ${application.containerAppName}        \
                        -g ${application.resourceGroups.instance} \
                        --revision ${revision.name}               \
                `)
            }
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
