const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const b2c     = require('./b2c')

async function main() {
    try {
        config.requiredEnvironmentVariableNonSecrets(['tag'])
        logging.printConfiguration(config)

        logging.header("Deploying teamsrv")

        const openIdConfig = await b2c.getWellKnownOpenIdConfiguration()

        const envVars = [
            `Login__OAuth__0__Module="login.jwt"`,
            `Login__OAuth__0__Options__JWKS__Uri="${openIdConfig.jwks_uri}"`,
            `Login__Basic__0__Module="login.pwd"`,
            `Login__Basic__0__Options__Credentials__0__Username="${config.authorizationServiceUsername}"`,
            `Login__Basic__0__Options__Credentials__0__Password=secretref:authorization-service-password`,
            `OpenApi__AuthorizationUrl="${openIdConfig.authorization_endpoint}"`,
            `OpenApi__TokenURL="${openIdConfig.token_endpoint}"`,
            `OpenApi__ClientId="${config.openApiClientId}"`,
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
                -n ${config.prefix}                           \
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

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
