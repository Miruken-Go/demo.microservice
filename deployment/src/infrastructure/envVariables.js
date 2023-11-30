import fs   from 'node:fs'
import path from 'node:path'


const env      = process.env.env
const instance = process.env.instance

export const variables = {
    env,
    instance,
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
    requireEnvFileVariables: function(directory, names){
        const filePath = path.join(directory, `${env}.json`)
        if (!fs.existsSync(filePath)) {
            throw new Error(`Config file does not exist: ${filePath}`)
        }

        const envSpecific = JSON.parse(fs.readFileSync(filePath))

        names.forEach(function(name) {
            const variable =  envSpecific[name]
            if (!variable){
                throw `Variable required from ${filePath}: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    }
}
