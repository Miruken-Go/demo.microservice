const {ApplicationType} = require('./infrastructure/systemDescription')

const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const systemDescription = {
    systemName: 'teamsrv',
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

const systemName   = systemDescription.systemName.toLowerCase() //must be lowercase
const appName      = systemDescription.appName.toLowerCase() //must be lowercase
const repository   = systemDescription.repository
const location     = systemDescription.location
const globalPrefix = `${appName}-global`
const commonPrefix = `${appName}-${env}`
const instance     = process.env.instance

const prefix = (instance) 
    ? `${appName}-${env}-${instance}`
    : `${appName}-${env}`

const b2cName                 = `${systemName}identity${env}`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
const b2cDisplayName          = `${systemName} identity ${env}`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
const b2cDomainName           = `${b2cName}.onmicrosoft.com`
const openIdConfigurationUrl  = `https://${b2cDisplayName}.b2clogin.com/${b2cDisplayName}.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN`
const keyVaultName            = `${commonPrefix}-keyvault` 
const containerRepositoryName = `${appName}global`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()

if (containerRepositoryName.length > 32)
    throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${containerRepositoryName} [${containerRepositoryName.length}]`

const imageName = `${containerRepositoryName}.azurecr.io/${appName}`

const config = {
    systemDescription,
    env,
    instance,
    appName,
    prefix,
    containerRepositoryName,
    imageName,
    keyVaultName,
    b2cName,       
    b2cDisplayName,
    b2cDomainName,
    openIdConfigurationUrl,
    location,
    repository,
    workingDirectory:                 process.cwd(),
    nodeDirectory:                    __dirname,
    defaultContainerImage:            'defaultContainerImage',
    globalResourceGroup:              globalPrefix,
    commonEnvironmentResourceGroup:   `${commonPrefix}-common`,
    manualEnvironmentResourceGroup:   `${commonPrefix}-manual`,
    environmentInstanceResourceGroup: `${prefix}`,
    secrets: {},
    requiredEnvironmentVariableSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name].trim()
            if (!variable){
                throw `Environment variable secret required: ${name}`
            }
            this.secrets[name] = variable
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

config.requiredEnvironmentVariableSecrets([
    'deploymentPipelineClientSecret',
])

config.requiredEnvironmentVariableNonSecrets([
    'tenantId',
    'subscriptionId',
    'deploymentPipelineClientId',
])

module.exports = {
    ...config,
}
