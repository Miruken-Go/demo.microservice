import * as az from './az.mjs'

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
