const config             = require('./config');
const keyvault           = require('./infrastructure/keyvault')
const b2cAppRegistration = require('./infrastructure/b2cAppRegistration') 
const {ApplicationType}  = require('./infrastructure/systemDescription')

async function main() {
    try {
        config.requiredEnvFileNonSecrets([
            'b2cDeploymentPipelineClientId',
        ])
        await keyvault.requireSecrets([
            'b2cDeploymentPipelineClientSecret',
        ])

        //Helpful documents about the application manifest
        //https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-1.0
        //https://learn.microsoft.com/en-us/graph/api/resources/webapplication?view=graph-rest-1.0
        //https://learn.microsoft.com/en-us/graph/api/resources/implicitgrantsettings?view=graph-rest-1.0
        //https://learn.microsoft.com/en-us/entra/identity-platform/reference-app-manifest

        const GROUPS_ID       = 'db580dbb-797c-4334-bf09-db802106accd'
        const ROLES_ID        = '0ba8756b-c67c-4fd3-9d70-488fc8da3b55'
        const ENTITLEMENTS_ID = 'd748b2c9-a76b-47b2-8c7b-fa348fbb474d'
      
        for (const app of config.systemDescription.applications) {
            const api = await b2cAppRegistration.createOrUpdateApplication({
                displayName:    app.name,
                signInAudience: 'AzureADandPersonalMicrosoftAccount',
                identifierUris: [ `https://${config.b2cName}.onmicrosoft.com/teamsrv` ],
                api: {
                    requestedAccessTokenVersion: 2,
                    oauth2PermissionScopes: [
                        {
                            id:                      GROUPS_ID, 
                            adminConsentDescription: 'Groups to which the user belongs.',
                            adminConsentDisplayName: 'Groups',
                            isEnabled:               true,
                            type:                    'User',
                            value:                   'Groups',
                        },
                        {
                            id:                      ROLES_ID, 
                            adminConsentDescription: 'Roles to which a user belongs',
                            adminConsentDisplayName: 'Roles',
                            isEnabled:               true,
                            type:                    'User',
                            value:                   'Roles',
                        },
                        {
                            id:                      ENTITLEMENTS_ID, 
                            adminConsentDescription: 'Entitlements which belong to the user',
                            adminConsentDisplayName: 'Entitlements',
                            isEnabled:               true,
                            type:                    'User',
                            value:                   'Entitlements',
                        },
                    ],
                }
            })
            console.log(api)
            
            const redirectUris = (['dev', 'qa'].includes(config.env))
                ? [
                    'https://jwt.ms',
                    'http://localhost:8080/oauth2-redirect.html',
                ] 
                : []

            const openapiUI = await b2cAppRegistration.createOrUpdateApplication({
                displayName:    `${app.name}_openapi`,
                signInAudience: 'AzureADandPersonalMicrosoftAccount',
                requiredResourceAccess: [
                    {
                        resourceAppId: api.appId, 
                        resourceAccess: [
                            {
                                id:   GROUPS_ID,
                                type: 'Scope',
                            },
                            {
                                id:   ROLES_ID,
                                type: 'Scope',
                            },
                            {
                                id:   ENTITLEMENTS_ID,
                                type: 'Scope',
                            },
                        ],
                    },
                ],
                web: {
                    implicitGrantSettings: {
                        enableAccessTokenIssuance: true,
                        enableIdTokenIssuance:     true,
                    }
                },
                spa: {
                    redirectUris: redirectUris
                },
            })
            console.log(openapiUI)
        }

        // const updatedApp = await b2cAppRegistration.addRedirectUris(openapiUI.id, ['https://foo.bar'])
        // console.log(updatedApp)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
