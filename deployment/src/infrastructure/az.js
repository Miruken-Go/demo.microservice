const bash       = require('./bash')
const config     = require('../config')
const { header } = require('./logging')

let loggedInToAZ  = false 
let loggedInToACR = false 

async function login() {
    if (loggedInToAZ) return 

    header('Logging into az')
    await bash.execute(`az login --service-principal --username ${config.deploymentPipelineClientId} --password ${config.secrets.deploymentPipelineClientSecret} --tenant ${config.tenantId}`);
    loggedInToAZ = true 
}

async function loginToACR() {
    if (loggedInToACR) return 

    header('Logging into ACR')
    await login()
    await bash.execute(`
        az acr login -n ${config.containerRepositoryName}
    `)
    loggedInToACR = true
}

async function createResourceGroup(name) {
    await login()
    await bash.execute(`az group create --location ${config.location} --name ${name} --subscription ${config.subscriptionId}`)
}

//https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-register-resource-provider?tabs=azure-cli
async function registerAzureProvider(providerName) { 
    await login()
    header(`Checking ${providerName} Provider Registration`)
    const providers = await bash.json(`az provider list --query "[?namespace=='${providerName}']" --output json`)
    if (providers.length) {
        const provider =  providers[0];
        if (provider.registrationState === "NotRegistered") {
            header(`Registering ${providerName} Provider`)
            await bash.execute(`az provider register --namespace ${providerName} --wait`);
        }
    }
}

async function getAzureContainerRepositoryPassword(name) {
    await login()
    const result = await bash.json(`az acr credential show --name ${name} --subscription ${config.subscriptionId}`, true)
    if (!result.passwords.length)
        throw new `Expected passwords from the Azure Container Registry: ${name}`

    return result.passwords[0].value
}

async function getKeyVaultSecret(secretName, keyVaultName) {
    await login()
    try {
        const result = await bash.json(`az keyvault secret show --name ${secretName} --vault-name ${keyVaultName}`, true)
        console.log(`Secret [${secretName}] found in [${keyVaultName}] keyvault`)
        return result.value
    } catch (error) {
        console.log(`Secret [${secretName}] not found in [${keyVaultName}] keyvault`)
        return null
    }
}

async function getContainerAppUrl(name) {
    await login()
    const result = await bash.json(`
        az containerapp show -n ${name} --resource-group ${config.environmentInstanceResourceGroup}
    `)
    return result.properties.configuration.ingress.fqdn
}

module.exports = {
    login,
    loginToACR,
    createResourceGroup,
    registerAzureProvider,
    getAzureContainerRepositoryPassword,
    getKeyVaultSecret,
    getContainerAppUrl,
}