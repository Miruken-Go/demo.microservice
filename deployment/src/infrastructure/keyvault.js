const az     = require('./az')
const config = require('../config')

const secrets = {}

async function requireSecrets() {
    for (const name of config.requiredKeyVaultSecrets) {
        if (secrets[name]) continue 

        const secret = await az.getKeyVaultSecret(name, config.keyVaultName)
        if (!secret){
            throw `KeyVault secret required: ${name}`
        }
        secrets[name] = secret.trim()
    }
}

module.exports = {
    requireSecrets,
    secrets
}
