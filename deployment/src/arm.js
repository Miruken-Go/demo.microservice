const { header } = require('./logging')
const path       = require('path')
const bash       = require('./bash')
const config     = require('./config')

async function deployGlobalResources() {
    header("Deploying Shared Arm Template")

    const bicepFile = path.join(config.nodeDirectory, 'bicep/globalResources.bicep')

    return await bash.json(`
        az deployment group create                                        \
            --template-file ${bicepFile}                                  \
            --subscription ${config.subscriptionId}                       \
            --resource-group ${config.globalResourceGroup}                \
            --mode complete                                               \
            --parameters                                                  \
                containerRepositoryName=${config.containerRepositoryName} \
                location=${config.location}                               \
    `)
}

async function deployEnvironmentInstanceResources(containerRepositoryPassword) {
    header("Deploying Environment Arm Template")

    const bicepFile = path.join(config.nodeDirectory, 'bicep/environmentInstanceResources.bicep')

    return await bash.json(`
        az deployment group create                                         \
            --template-file ${bicepFile}                                   \
            --subscription ${config.subscriptionId}                        \
            --resource-group ${config.resourceGroup}                       \
            --mode complete                                                \
            --parameters                                                   \
                prefix=${config.prefix}                                    \
                appName=${config.appName}                                  \
                location=${config.location}                                \
                containerRepositoryName=${config.containerRepositoryName}  \
                containerRepositoryPassword=${containerRepositoryPassword} \
    `)
}

module.exports = {
    deployEnvironmentInstanceResources,
    deployGlobalResources,
}
