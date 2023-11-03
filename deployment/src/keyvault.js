const az     = require('./az')
const config = require('./config')

const secrets = {}

async function requireSecrets() {
    const keyVaultName = `teamsrv-pipeline-${config.env}` 
    for (const name of config.requiredKeyVaultSecrets) {
        if (secrets[name]) continue 

        const secret = await az.getKeyVaultSecret(name, keyVaultName)
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
