const az      = require('./az');
const arm     = require('./arm');
const logging = require('./logging');
const config  = require('./config');

async function main() {
    try {
        console.log("Deploying Environment")
        logging.printConfiguration(config)
        await az.login()

        //Environment resources
        await az.createResourceGroup(config.resourceGroup)

        const getAzureContainerRepositoryPassword = await az.getAzureContainerRepositoryPassword(config.containerRepositoryName)
        await arm.deployEnvironmentResources(getAzureContainerRepositoryPassword)

        console.log("Configuration Succeded")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Configuration Failed")
    }
}

main()
