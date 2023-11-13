const { Organization } = require('./infrastructure/config')

const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const instance = process.env.instance

const org = new Organization({
    name:     'MajorLeageMiruken',
    location: 'CentralUs',
    env:      env,
    instance: instance,
    domains: [
        {
            name: 'billing', 
            apps: [
                {
                    name: 'billingui',  
                    ui:   true
                },
                {
                    name: 'billingsrv', 
                    ui:   true, 
                    api:  true
                },
            ]
        },
        {
            name: 'league', 
            apps: [
                {
                    name: 'majorleaguemiruken', 
                    ui:   true
                },
                {
                    name: 'tournaments',
                    ui:   true
                },
                {
                    name: 'teamsrv',            
                    ui:   true, 
                    api:  true
                },
                {
                    name: 'schedulesrv',        
                    ui:   true, 
                    api:  true
                },
            ]
        },
    ],
})



const systemDescription = {
    orgName: 'teamsrv',
    applications: [
        {
            name: 'teamsrv',
            type: ApplicationType.apiWithOpenApiUI
        }
    ],
    appName:    'teamsrv',
    repository: 'https://github.com/Miruken-Go/demo.microservice',
    location:   'CentralUs',
    environments: [
        'dev',
        'qa',
        'uat',
        'demo',
        'prod',
        'dr'
    ]
}

//const orgName      = systemDescription.orgName.toLowerCase() //must be lowercase
//const appName      = systemDescription.appName.toLowerCase() //must be lowercase
//const repository   = systemDescription.repository  // Should go away. Should be part of the build.
//const location     = systemDescription.location
//const globalPrefix = `${appName}-global`
//const commonPrefix = `${appName}-${env}`

// const prefix = (instance) 
//     ? `${appName}-${env}-${instance}`
//     : `${appName}-${env}`

//const b2cName                 = `${orgName}identity${env}`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
//const b2cDisplayName          = `${orgName} identity ${env}`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
//const b2cDomainName           = `${b2cName}.onmicrosoft.com`
//const openIdConfigurationUrl  = `https://${b2cDisplayName}.b2clogin.com/${b2cDisplayName}.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN`
//const keyVaultName            = `${commonPrefix}-keyvault` 
/*
    const containerRepositoryName = `${appName}global`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
    if (containerRepositoryName.length > 32)
        throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${containerRepositoryName} [${containerRepositoryName.length}]`
*/


//Needs an app and an organization to build this
//const imageName = `${containerRepositoryName}.azurecr.io/${appName}` 

const config = {
    systemDescription,
    env,
    instance,
    workingDirectory:                 process.cwd(),
    nodeDirectory:                    __dirname,
    defaultContainerImage:            'defaultContainerImage',
    secrets: {},

    requiredEnvironmentVariableSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name]
            if (!variable){
                throw `Environment variable secret required: ${name}`
            }
            this.secrets[name] = variable.trim()
        }.bind(this));
    },
    requiredEnvironmentVariableNonSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name] || this[name]
            if (!variable){
                throw `Environment variable required: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    }, 
    requiredEnvFileNonSecrets: function(names){
        if (systemDescription.environments.includes(env)) {
            const envSpecific = require(`./${env}.js`)
            names.forEach(function(name) {
                const variable =  envSpecific[name]
                if (!variable){
                    throw `Variable required from ${env}.js: ${name}`
                }
                this[name] = variable.trim()
            }.bind(this));
        }
    }
}


// config.requiredEnvironmentVariableSecrets([
//     'deploymentPipelineClientSecret',
// ])

// config.requiredEnvironmentVariableNonSecrets([
//     'tenantId',
//     'subscriptionId',
//     'deploymentPipelineClientId',
// ])

module.exports = {
    ...config
}
