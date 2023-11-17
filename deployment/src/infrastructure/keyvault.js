const az = require('./az')

const secrets = {
    async requireSecrets (names, keyVaultName) {
        names.forEach(async function(name) {
            if (this[name]) return 

            const secret = await az.getKeyVaultSecret(name, keyVaultName)
            if (!secret){
                throw `KeyVault secret required: ${name}`
            }
            this[name] = secret.trim()
        }.bind(this));
    }
}

module.exports = {
    secrets
}
