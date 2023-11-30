import * as az from './az.js'

export const secrets = {
    async requireSecrets (names, keyVaultName) {
        for(const name of names) {
            if (this[name]) return 

            const secret = await az.getKeyVaultSecret(name, keyVaultName)
            if (!secret){
                throw `KeyVault secret required: ${name}`
            }
            this[name] = secret.trim()
        }
    }
}
