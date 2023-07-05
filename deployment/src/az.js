const bash       = require('./bash')
const config     = require('./config')
const { header } = require('./logging')

async function login() {
    header('Logging into az')
    await bash.execute(`az login --service-principal --username ${config.deploymentPipelineClientId} --password ${config.secrets.deploymentPipelineClientSecret} --tenant ${config.tenantId}`);
}

async function createResourceGroup(name) {
    await bash.execute(`az group create --location ${config.location} --name ${name} --subscription ${config.subscriptionId}`)
}

//https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-register-resource-provider?tabs=azure-cli
async function registerAzureProvider(providerName) { 
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
    const result = await bash.json(`az acr credential show --name ${name} --subscription ${config.subscriptionId}`, true)
    if (!result.passwords.length)
        throw new `Expected passwords from the Azure Container Registry: ${name}`

    return result.passwords[0].value
}

module.exports = {
    login,
    createResourceGroup,
    registerAzureProvider,
    getAzureContainerRepositoryPassword
}