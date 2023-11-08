const az     = require('./az')
const config = require('../config')

const secrets = {}

async function requireSecrets(names) {
    for (const name of names) {
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
