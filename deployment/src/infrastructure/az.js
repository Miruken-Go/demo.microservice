import * as bash     from './bash.js'
import { header }    from './logging.js'
import { variables } from './envVariables.js'
import { secrets }   from './envSecrets.js'

variables.requireEnvVariables([
    'tenantId',
    'subscriptionId',
    'deploymentPipelineClientId',
])

secrets.require([
    'deploymentPipelineClientSecret',
])

let loggedInToAZ  = false 
let loggedInToACR = false 

export async function login() {
    if (loggedInToAZ) return 

    header('Logging into az')
    await bash.execute(`az login --service-principal --username ${variables.deploymentPipelineClientId} --password ${secrets.deploymentPipelineClientSecret} --tenant ${variables.tenantId}`);
    loggedInToAZ = true 
}

export async function loginToACR(containerRepositoryName) {
    if (loggedInToACR) return 

    header('Logging into ACR')
    await login()
    await bash.execute(`
        az acr login -n ${containerRepositoryName}
    `)
    loggedInToACR = true
}

export async function createResourceGroup(name, location, tags) {
    await login()
    await bash.execute(`az group create --location ${location} --name ${name} --subscription ${variables.subscriptionId} --tags ${tags}`)
}

//https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-register-resource-provider?tabs=azure-cli
export async function registerAzureProvider(providerName) { 
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

export async function getAzureContainerRepositoryPassword(name) {
    await login()
    const result = await bash.json(`az acr credential show --name ${name} --subscription ${variables.subscriptionId}`, true)
    if (!result.passwords.length)
        throw new `Expected passwords from the Azure Container Registry: ${name}`

    return result.passwords[0].value
}

export async function getKeyVaultSecret(secretName, keyVaultName) {
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

export async function getContainerAppUrl(name, resourceGroup) {
    await login()
    const result = await bash.json(`
        az containerapp show -n ${name} --resource-group ${resourceGroup}
    `)

    if (!result) throw new Error(`ContainerApp ${name} not found in ${resourceGroup}`)
    
    return result.properties.configuration.ingress.fqdn
}


export async function deleteOrphanedApplicationSecurityPrincipals(name) {
    await login()
    const ids = await bash.json(`
        az role assignment list --all --query "[?principalName==''].id"    
    `)

    if (ids.length) {
        await bash.json(`
            az role assignment delete --ids ${ids.join(' ')}
        `)

        console.log(`Deleted ${ids.length} orphaned application security principals`)
    }
}
