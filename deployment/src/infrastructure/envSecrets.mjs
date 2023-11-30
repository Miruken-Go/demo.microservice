export const secrets = {
    require: function (names) {
        names.forEach(function(name) {
            if(this[name]) return

            const secret = process.env[name]
            if (!secret){
                throw `Environment secret required: ${name}`
            }
            this[name] = secret.trim()
        }.bind(this));
    }
}
