const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const appDescription ={
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

const appName      = appDescription.appName.toLowerCase() //must be lowercase
const repository   = appDescription.repository
const location     = appDescription.location
const globalPrefix = `${appName}-global`
const commonPrefix = `${appName}-${env}`
const instance     = process.env.instance

const prefix = (instance) 
    ? `${appName}-${env}-${instance}`
    : `${appName}-${env}`

const keyVaultName            = `${commonPrefix}-keyvault` 
const containerRepositoryName = `${appName}global`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()

if (containerRepositoryName.length > 32)
    throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${containerRepositoryName} [${containerRepositoryName.length}]`

const imageName = `${containerRepositoryName}.azurecr.io/${appName}`

const config = {
    env,
    instance,
    appName,
    prefix,
    containerRepositoryName,
    imageName,
    keyVaultName,
    location,
    repository,
    workingDirectory:                 process.cwd(),
    nodeDirectory:                    __dirname,
    defaultContainerImage:            'defaultContainerImage',
    globalResourceGroup:              globalPrefix,
    commonEnvironmentResourceGroup:   `${commonPrefix}-common`,
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
        if (appDescription.environments.includes(env)) {
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

config.requiredEnvFileNonSecrets([
    'b2cDeploymentPipelineClientId',
    'identityExperienceFrameworkClientId',
    'proxyIdentityExperienceFrameworkClientId',
    'openapiClientId',
    'b2cDomainName',
    'wellKnownOpenIdConfigurationUrl',
    'authorizationServiceUrl',
    'authorizationServiceUsername',
])

config.requiredKeyVaultSecrets = [
    'b2cDeploymentPipelineClientSecret',
]

module.exports = {
    ...config,
}
