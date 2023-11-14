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

/*
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
*/


// config.requiredEnvironmentVariableSecrets([
//     'deploymentPipelineClientSecret',
// ])

// config.requiredEnvironmentVariableNonSecrets([
//     'tenantId',
//     'subscriptionId',
//     'deploymentPipelineClientId',
// ])

module.exports = {
    workingDirectory: process.cwd(),
    nodeDirectory:    __dirname,
    organization: org
}
