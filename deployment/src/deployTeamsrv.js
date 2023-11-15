const config             = require('./config');
const logging            = require('./infrastructure/logging');
const az                 = require('./infrastructure/az');
const bash               = require('./infrastructure/bash')
const b2c                = require('./infrastructure/b2c')
const keyvault           = require('./infrastructure/keyvault')
const b2cAppRegistration = require('./infrastructure/b2cAppRegistration')

async function main() {
    try {
        config.requiredEnvironmentVariableNonSecrets(['tag'])
        config.requiredEnvFileNonSecrets([
            'b2cDeploymentPipelineClientId',
            'authorizationServiceUsername',
        ])
        await keyvault.requireSecrets([
            'b2cDeploymentPipelineClientSecret',
        ])
        logging.printConfiguration(config)

        const applicationName = config.prefix

        logging.header(`Deploying ${applicationName}`)

        const openIdConfig    = await b2c.getWellKnownOpenIdConfiguration()
        const teamsrv         = await b2cAppRegistration.getApplicationByName('teamsrv')
        const teamsrv_openapi = await b2cAppRegistration.getApplicationByName('teamsrv_openapi')

        const envVars = [
            `Login__OAuth__0__Module="login.jwt"`,
            `Login__OAuth__0__Options__Audience="${teamsrv.appId}"`,
            `Login__OAuth__0__Options__JWKS__Uri="${openIdConfig.jwks_uri}"`,
            `Login__Basic__0__Module="login.pwd"`,
            `Login__Basic__0__Options__Credentials__0__Username="${config.authorizationServiceUsername}"`,
            `Login__Basic__0__Options__Credentials__0__Password=secretref:authorization-service-password`,
            `OpenApi__AuthorizationUrl="${openIdConfig.authorization_endpoint}"`,
            `OpenApi__TokenUrl="${openIdConfig.token_endpoint}"`,
            `OpenApi__ClientId="${teamsrv_openapi.appId}"`,
            `OpenApi__OpenIdConnectUrl="${config.wellKnownOpenIdConfigurationUrl}"`,
            `OpenApi__Scopes__0__Name="https://${config.b2cDomainName}/teamsrv/Groups"`,
            `OpenApi__Scopes__0__Description="Groups to which the user belongs."`,
            `OpenApi__Scopes__1__Name="https://${config.b2cDomainName}/teamsrv/Roles"`,
            `OpenApi__Scopes__1__Description="Roles to which the user belongs."`,
            `OpenApi__Scopes__2__Name="https://${config.b2cDomainName}/teamsrv/Entitlements"`, 
            `OpenApi__Scopes__2__Description="Entitlements the user has."`,
        ]

        await az.login()

        //https://learn.microsoft.com/en-us/cli/azure/containerapp?view=azure-cli-latest#az-containerapp-update
        //Create the new revision
        await bash.execute(`
            az containerapp update                            \
                -n ${applicationName}                         \
                -g ${config.environmentInstanceResourceGroup} \
                --image ${config.imageName}:${config.tag}     \
                --container-name ${config.appName}            \
                --revision-suffix ${config.tag}               \
                --replace-env-vars ${envVars.join(' ')}       \
        `)

        const revisions = await bash.json(`
            az containerapp revision list  \
                -n ${config.prefix}        \
                -g ${config.environmentInstanceResourceGroup} \
        `)

        const revisionsToActivate   = revisions.filter(r => r.name.includes(config.tag))
        const revisionsToDeactivate = revisions.filter(r => !r.name.includes(config.tag))

        //You must have an active revision before deactivating the rest
        for (const revision of revisionsToActivate) {
            if (revision.properties.active !== true) {
                console.log(`Activating: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision activate \
                        -n ${config.prefix}           \
                        -g ${config.environmentInstanceResourceGroup}    \
                        --revision ${revision.name}   \
                `)
            }

            await bash.execute(`
                az containerapp ingress traffic set \
                    -n ${config.prefix}           \
                    -g ${config.environmentInstanceResourceGroup}    \
                    --revision-weight ${revision.name}=100
            `)
        }

        for (const revision of revisionsToDeactivate) {
            if (revision.properties.active === true) {
                console.log(`Dectivate: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision deactivate \
                        -n ${config.prefix}           \
                        -g ${config.environmentInstanceResourceGroup}    \
                        --revision ${revision.name}   \
                `)
            }
        }

        const appUrl = await az.getContainerAppUrl(applicationName)
        await b2cAppRegistration.addRedirectUris(teamsrv_openapi.id, [`https://${appUrl}`])

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
