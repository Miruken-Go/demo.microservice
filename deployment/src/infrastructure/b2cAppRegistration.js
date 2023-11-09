const config = require('../config');
const graph  = require('./graph');
const az     = require('./az');

async function getApplications() {
    const result = await graph.get("/applications")
    return result.data.value
}

async function getApplicationById(id) {
    const result = await graph.get(`/applications/${id}`)
    return result.data
}

async function getApplicationByName(displayName) {
    const applications = await getApplications()
    const application  =  applications.find(a => a.displayName === displayName)
    console.log(application)
    return application
}

async function updateApplication(id, manifest) {
    console.log(`Updating existing application appId [${id}]`)
    await graph.patch(`/applications/${id}`, manifest)
    return await getApplicationById(id)
}

async function createOrUpdateApplication(manifest) {
    const displayName = manifest.displayName
    const existing = await getApplicationByName(displayName)

    let application = undefined
    if (existing) {
        application = await updateApplication(existing.id, manifest)
    } else {
        console.log(`Creating application: ${displayName}`)
        application = (await graph.post("/applications", manifest)).data
        console.log(`Created Application: ${displayName}`)
        console.log(application)
    }

    return application
}

async function addRedirectUris(id, uris) {
    const app = await getApplicationById(id)
    const redirectUris = [...app.spa.redirectUris]
    for (const uri of uris) {
        if (!redirectUris.includes(uri)) {
            redirectUris.push(uri)
        }
    }
    await updateApplication(id, {
        spa: {
            redirectUris: redirectUris
        }
    })
    return await getApplicationById(id)
}

async function configure() {
    const GROUPS_ID       = 'db580dbb-797c-4334-bf09-db802106accd'
    const ROLES_ID        = '0ba8756b-c67c-4fd3-9d70-488fc8da3b55'
    const ENTITLEMENTS_ID = 'd748b2c9-a76b-47b2-8c7b-fa348fbb474d'
    
    for (const app of config.systemDescription.applications) {
        const api = await createOrUpdateApplication({
            displayName:    app.name,
            signInAudience: 'AzureADandPersonalMicrosoftAccount',
            identifierUris: [ `https://${config.b2cName}.onmicrosoft.com/${app.name}` ],
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

        const appUrl = await az.getContainerAppUrl(config.prefix)
        if(!appUrl) throw new Error(`default application redirectUri could not be calculated. The AppUrl for ${config.prefix} container app was not found. The default application environment instance needs to be deployed before common configuration can run.`)
        const appRedirectUri = `https://${appUrl}`
        
        const redirectUris = (['dev', 'qa'].includes(config.env))
            ? [
                appRedirectUri,
                'https://jwt.ms',
                'http://localhost:8080/oauth2-redirect.html',
              ] 
            : [

                appRedirectUri,
              ]

        const openapiUI = await createOrUpdateApplication({
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
}

module.exports = {
    getApplications,
    getApplicationById,
    getApplicationByName,
    updateApplication,
    createOrUpdateApplication,
    addRedirectUris,
    configure,
}
