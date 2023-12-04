import * as logging  from '#infrastructure/logging.js'
import * as az       from '#infrastructure/az.js'
import * as bash     from '#infrastructure/bash.js'
import * as gh       from '#infrastructure/gh.js'
import { B2C }       from '#infrastructure/b2c.js'
import { variables } from '#infrastructure/envVariables.js'

import { 
    configDirectory,
    organization 
} from './config.js'

variables.requireEnvVariables([
    'tag'
])

variables.requireEnvFileVariables(configDirectory, [
    'b2cDeploymentPipelineClientId'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        const application = organization.getApplicationByName("team-srv")

        logging.header(`Deploying ${application.name}`)

        const b2c             = new B2C(organization, variables.b2cDeploymentPipelineClientId)
        const appRegistration = await b2c.getApplicationByName(application.parent.name)
        const openIdConfig    = await b2c.getWellKnownOpenIdConfiguration()

        const envVars = [
            `Login__OAuth__0__Module="login.jwt"`,
            `Login__OAuth__0__Options__Audience="${appRegistration.appId}"`,
            `Login__OAuth__0__Options__JWKS__Uri="${openIdConfig.jwks_uri}"`,
            `OpenApi__AuthorizationUrl="${openIdConfig.authorization_endpoint}"`,
            `OpenApi__TokenUrl="${openIdConfig.token_endpoint}"`,
            `OpenApi__ClientId="${appRegistration.appId}"`,
            `OpenApi__OpenIdConnectUrl="${organization.b2c.openIdConfigurationUrl}"`,
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
        //We have to add the current time to the revision suffix to make it unique
        //otherwise we would never be able to redeploy the same container tag
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

        const appUrl = await az.getContainerAppUrl(application.containerAppName, application.resourceGroups.instance)
        await b2c.addRedirectUris(appRegistration.id, [`https://${appUrl}/oauth2-redirect.html`])

        await gh.sendRepositoryDispatch(`deployed-${application.name}`, {
            tag: variables.tag
        })

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
