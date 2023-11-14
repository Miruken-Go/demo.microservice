const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const instance = process.env.instance

const variables = {
    env,
    instance,
    require: function(names) {
        names.forEach(function(name) {
            if(this[name]) return

            const variable = process.env[name]
            if (!variable){
                throw `Environment variable required: ${name}`
            }
            this[name] = variable.trim()
        }.bind(this));
    }
}

module.exports = {
    variables
}
