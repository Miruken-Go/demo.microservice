const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const appName         = 'teamsrv'.toLowerCase() //must be lowercase
const repository      = 'https://github.com/Miruken-Go/demo.microservice'
const defaultLocation = 'CentralUs'

const instance = process.env.instance
const prefix   = (instance) 
    ? `${appName}-${env}-${instance}`
    : `${appName}-${env}`

const simplePrefix = prefix.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
const containerRepositoryName = `${appName}shared`.replace(/[^A-Za-z0-9]/g, "").toLowerCase()

if (containerRepositoryName.length > 32)
    throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${containerRepositoryName} [${containerRepositoryName.length}]`

if (simplePrefix.length > 32)
    throw `Configuration Error - simplePrefix cannot be longer than 50 characters because of ACA naming restrictions: ${simplePrefix} [${simplePrefix.length}]`

const imageName = `${containerRepositoryName}.azurecr.io/${appName}`

const envSpecific = require(`./${env}.js`)

const config = {
    workingDirectory: process.cwd(),
    nodeDirectory:    __dirname,
    env,
    instance,
    appName,
    defaultContainerImage: 'defaultContainerImage',
    prefix,
    simplePrefix,
    resourceGroup: `${prefix}-rg`,
    globalResourceGroup: `${appName}-global`,
    containerRepositoryName,
    imageName,
    location: process.env.location || defaultLocation,
    repository,
    secrets: {},
    ...envSpecific,
    requiredSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name].trim()
            if (!variable){
                throw `Environment variable secret required: ${name}`
            }
            this.secrets[name] = variable
        }.bind(this));
    },
    requiredNonSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name] || this[name]
            if (!variable){
                throw `Environment variable required: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    }, 
}

config.requiredSecrets([
    'deploymentPipelineClientSecret'
])

config.requiredNonSecrets([
    'tenantId',
    'subscriptionId',
    'deploymentPipelineClientId',
    'b2cDeploymentPipelineClientId',
    'identityExperienceFrameworkClientId',
    'proxyIdentityExperienceFrameworkClientId',
    'b2cDomainName',
    'authorizatioServiceUrl',
])

config.requiredKeyVaultSecrets = [
    'b2cDeploymentPipelineClientSecret',
]

module.exports = {
    ...config,
}
