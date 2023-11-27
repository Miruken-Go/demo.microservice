const logging            = require('./infrastructure/logging');
const az                 = require('./infrastructure/az');
const bash               = require('./infrastructure/bash')
const { B2C }            = require('./infrastructure/b2c')
const { variables }      = require('./infrastructure/envVariables')
const { organization }   = require('./config');

variables.requireEnvVariables([
    'tag'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        const application = organization.getApplicationByName("adb2c-auth-srv")

        logging.header(`Deploying ${application.name}`)

        const b2c             = new B2C(organization)
        const appRegistration = await b2c.getApplicationByName(application.parent.name)
        const openIdConfig    = await b2c.getWellKnownOpenIdConfiguration()

        const envVars = [
            `Login__OAuth__0__Module='login.jwt'`,
            `Login__OAuth__0__Options__Audience='${appRegistration.appId}'`,
            `Login__OAuth__0__Options__JWKS__Uri='${openIdConfig.jwks_uri}'`,
            `OpenApi__AuthorizationUrl='${openIdConfig.authorization_endpoint}'`,
            `OpenApi__TokenUrl='${openIdConfig.token_endpoint}'`,
            `OpenApi__ClientId='${appRegistration.appId}'`,
            `OpenApi__OpenIdConnectUrl='${organization.b2c.openIdConfigurationUrl}'`,
        ]

        const identifierUri = appRegistration.identifierUris[0]
        const scopes        = appRegistration.api.oauth2PermissionScopes
        for (let i = 0; i < scopes.length; i++) {
            const scope = scopes[i]
            envVars.push(`OpenApi__Scopes__${i}__Name='${identifierUri}/${scope.value}'`)
            envVars.push(`OpenApi__Scopes__${i}__Description='${scope.adminConsentDescription}'`)
        }

        await az.login()

        //https://learn.microsoft.com/en-us/cli/azure/containerapp?view=azure-cli-latest#az-containerapp-update
        //Create the new revision
        const now = `${Math.floor(Date.now()/1000)}`.trim()
        await bash.execute(`
            az containerapp update                                \
                -n ${application.containerAppName}                \
                -g ${application.resourceGroups.instance}         \
                --image ${application.imageName}:${variables.tag} \
                --container-name ${application.name}              \
                --revision-suffix ${variables.tag}-${now}         \
                --replace-env-vars ${envVars.join(' ')}           \
        `)

        const revisions = await bash.json(`
            az containerapp revision list                 \
                -n ${application.containerAppName}        \
                -g ${application.resourceGroups.instance} \
        `)

        const revisionsToActivate   = revisions.filter(r => r.name.includes(variables.tag))
        const revisionsToDeactivate = revisions.filter(r => !r.name.includes(variables.tag))

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

        const appUrl = await az.getContainerAppUrl(application.containerAppName, application.resourceGroups.instance)
        await b2c.addRedirectUris(appRegistration.id, [`https://${appUrl}/oauth2-redirect.html`])

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
