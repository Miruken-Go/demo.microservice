const env = process.env.env
//if (!env) throw "Environment variable required: [env]"

const instance = process.env.instance

const variables = {
    env,
    instance,
    workingDirectory: process.cwd(),
    nodeDirectory:    __dirname,
    optionalEnvVariables: function(names) {
        names.forEach(function(name) {
            if(this[name]) return

            const variable = process.env[name]
            if (variable){
                this[name] = variable.trim()
            }
        }.bind(this));
    },
    requireEnvVariables: function(names) {
        names.forEach(function(name) {
            if(this[name]) return

            const variable = process.env[name]
            if (!variable){
                throw `Environment variable required: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    },
    requireEnvFileVariables: function(names){
        const envSpecific = require(`../${env}.js`)
        names.forEach(function(name) {
            const variable =  envSpecific[name]
            if (!variable){
                throw `Variable required from ${env}.js: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    }
}

module.exports = {
    variables
}
