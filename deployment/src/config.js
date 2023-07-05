const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const appName  = 'teamsrv'.toLowerCase()
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

const config = {
    workingDirectory: process.cwd(),
    nodeDirectory:    __dirname,
    env,
    instance,
    appName,
    prefix,
    simplePrefix,
    resourceGroup: `${prefix}-rg`,
    sharedResourceGroup: `${appName}-shared-rg`,
    containerRepositoryName,
    imageName,
    location:      process.env.location || 'CentralUs',
    secrets:       {},

    requiredSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name]
            if (!variable){
                throw `Environment variable secret required: ${name}`
            }
            this.secrets[name] = variable
        }.bind(this));
    },
    requiredNonSecrets: function (names) {
        names.forEach(function(name) {
            const variable = process.env[name]
            if (!variable){
                throw `Environment variable required: ${name}`
            }
            this[name] = variable
        }.bind(this));
    }
}

config.requiredSecrets([
    'deploymentPipelineClientSecret'
])

config.requiredNonSecrets([
    'tenantId',
    'subscriptionId',
    'deploymentPipelineClientId'
])

//const envSpecific = require(`./${env}.js`)

module.exports = {
    ...config,
    //...envSpecific
}
