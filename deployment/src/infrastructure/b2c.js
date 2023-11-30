import * as logging  from './logging.js'
import * as az       from './az.js'
import { Graph }     from './graph.js'
import { variables } from './envVariables.js'
import * as fs       from 'node:fs'
import * as path     from 'node:path'
import axios         from 'axios'

variables.requireEnvVariables([
    'env',
])

export class B2C {
    organization
    graph

    constructor (organization, b2cDeploymentPipelineClientId) {
        if (!organization)                  throw new Error('organization is required')
        if (!b2cDeploymentPipelineClientId) throw new Error('b2cDeploymentPipelineClientId is required')

        this.organization = organization
        this.graph        = new Graph(organization, b2cDeploymentPipelineClientId)
    }

    async getWellKnownOpenIdConfiguration() {
        const uri = this.organization.b2c.openIdConfigurationUrl
        console.log(`Getting: ${uri}`)
        const result = await axios.get(uri)
            .catch(function (error) {
                console.log(`Failed to Get: ${uri}`)  
                logError(error)
                throw error
            });
    
        console.log(result.data)
        return result.data
    }

    async getApplications() {
        const result = await this.graph.get("/applications")
        return result.data.value
    }

    async getApplicationById(id) {
        const result = await this.graph.get(`/applications/${id}`)
        return result.data
    }

    async getApplicationByName(displayName) {
        const applications = await this.getApplications()
        const application  =  applications.find(a => a.displayName === displayName)
        console.log(application)
        return application
    }

    async updateApplication(id, manifest) {
        console.log(`Updating existing application appId [${id}]`)
        await this.graph.patch(`/applications/${id}`, manifest)
        return await this.getApplicationById(id)
    }

    async createOrUpdateApplication(manifest) {
        const displayName = manifest.displayName
        const existing = await this.getApplicationByName(displayName)

        let application = undefined
        if (existing) {
            application = await this.updateApplication(existing.id, manifest)
        } else {
            console.log(`Creating application: ${displayName}`)
            application = (await this.graph.post("/applications", manifest)).data
            console.log(`Created Application: ${displayName}`)
            console.log(application)
        }

        return application
    }

    async addRedirectUris(id, uris) {
        const app = await this.getApplicationById(id)
        const redirectUris = [...app.spa.redirectUris]
        for (const uri of uris) {
            if (!redirectUris.includes(uri)) {
                redirectUris.push(uri)
            }
        }
        await this.updateApplication(id, {
            spa: {
                redirectUris: redirectUris
            }
        })
        return await this.getApplicationById(id)
    }

    async configureAppRegistrations() { 
        await this.configureAppRegistration(this.organization)
        for (const domain of this.organization.domains) {
            await this.configureAppRegistration(domain)
        }
    }

    async configureAppRegistration(domain) {
        const GROUPS_ID       = 'db580dbb-797c-4334-bf09-db802106accd'
        const ROLES_ID        = '0ba8756b-c67c-4fd3-9d70-488fc8da3b55'
        const ENTITLEMENTS_ID = 'd748b2c9-a76b-47b2-8c7b-fa348fbb474d'
        
        const appRegistration = await this.createOrUpdateApplication({
            displayName:    domain.name,
            signInAudience: 'AzureADandPersonalMicrosoftAccount',
            identifierUris: [ `https://${this.organization.b2c.name}.onmicrosoft.com/${domain.name}` ],
            api: {
                requestedAccessTokenVersion: 2,
                oauth2PermissionScopes: [
                    {
                        id:                      GROUPS_ID, 
                        adminConsentDescription: 'Groups to which the user belongs.',
                        adminConsentDisplayName: 'Groups',
                        isEnabled:               true,
                        type:                    'Admin',
                        value:                   'Groups',
                    },
                    {
                        id:                      ROLES_ID, 
                        adminConsentDescription: 'Roles to which a user belongs',
                        adminConsentDisplayName: 'Roles',
                        isEnabled:               true,
                        type:                    'Admin',
                        value:                   'Roles',
                    },
                    {
                        id:                      ENTITLEMENTS_ID, 
                        adminConsentDescription: 'Entitlements which belong to the user',
                        adminConsentDisplayName: 'Entitlements',
                        isEnabled:               true,
                        type:                    'Admin',
                        value:                   'Entitlements',
                    },
                ],
            },
        })
        console.log(appRegistration)

        const requiredResourceAccess = await this.createOrUpdateApplication({
            displayName: domain.name,
            requiredResourceAccess: [
                {
                    resourceAppId: '00000003-0000-0000-c000-000000000000',
                    resourceAccess: [
                        {
                            'id':   '37f7f235-527c-4136-accd-4a02d197296e',
                            'type': 'Scope'
                        },
                        {
                            'id':   '7427e0e9-2fba-42fe-b0c0-848c9e6a8182',
                            'type': 'Scope'
                        }
                    ]
                },
                {
                    resourceAppId: appRegistration.appId, 
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
        })
        console.log(requiredResourceAccess)

        if (domain.applications.some(a => a.implicitFlow)) {
            console.log('Configure implicit flow')
            const implicicteGrant = await this.createOrUpdateApplication({
                displayName: domain.name,
                web: {
                    implicitGrantSettings: {
                        enableAccessTokenIssuance: true,
                        enableIdTokenIssuance:     true,
                    }
                }
            })
            console.log(implicicteGrant)
        }

        if (domain.applications.some(a => a.spa)) {
            console.log('Configure spa')
            const redirectUris = (['dev', 'qa'].includes(variables.env))
                ? [
                    'https://jwt.ms/',
                    'http://localhost:8080/oauth2-redirect.html',
                ] 
                : []
            const spa = await this.createOrUpdateApplication({
                displayName: domain.name,
                spa: {
                    redirectUris: redirectUris
                },
            })
            console.log(spa)
        }
    }

    async configureCustomPolicies () {

        logging.header("Deploying B2C Configuration")

        const identityExperienceFrameworkClient = await this.getApplicationByName('IdentityExperienceFramework')
        if(!identityExperienceFrameworkClient) throw new Error("IdentityExperienceFramework application not found. Check that the App Registration was created in B2C and check the name spelling and casing.")

        const proxyIdentityExperienceFrameworkClient = await this.getApplicationByName('ProxyIdentityExperienceFramework')
        if(!proxyIdentityExperienceFrameworkClient) throw new Error("ProxyIdentityExperienceFramework application not found. Check that the App Registration was created in B2C and check the name spelling and casing.")

        const containerAppName = this.organization.enrichApiApplication.containerAppName
        const appUrl = await az.getContainerAppUrl(containerAppName, this.organization.resourceGroups.instance)
        if(!appUrl) throw new Error(`authorizationServiceUrl could not be calculated. The AppUrl for ${containerAppName} container app was not found. The default application environment instance needs to be deployed before common configuration can run.`)

        const authorizationServiceUrl = `https://${appUrl}/enrich`

        //https://learn.microsoft.com/en-us/azure/active-directory-b2c/deploy-custom-policies-devops
        const customPoliciesDirectory = new URL('../custom-policies', import.meta.url).pathname
        const customPoliciesFileOrder = [
            'TrustFrameworkBase.xml',
            'TrustFrameworkLocalization.xml',
            'TrustFrameworkExtensions.xml',
            'SignUp_SignIn.xml',
            'ProfileEdit.xml',
            'PasswordReset.xml',
        ]
        for (const file of customPoliciesFileOrder) {
                const policyId = `B2C_1A_${path.basename(file, '.xml')}`
                const filePath = path.join(customPoliciesDirectory, file)
                let xml        = fs.readFileSync(filePath,{encoding: 'utf-8'}) 

                xml = xml.replace(/{B2C_DOMAIN_NAME}/g,                              this.organization.b2c.domainName)
                xml = xml.replace(/{IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g,       identityExperienceFrameworkClient.appId)
                xml = xml.replace(/{PROXY_IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g, proxyIdentityExperienceFrameworkClient.appId)
                xml = xml.replace(/{AUTHORIZATION_SERVICE_URL}/g,                    authorizationServiceUrl)

                await this.graph.updateTrustFrameworkPolicy(policyId, xml)
        };
    }
}
