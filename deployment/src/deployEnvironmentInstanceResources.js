const az       = require('./infrastructure/az');
const arm      = require('./infrastructure/arm');
const logging  = require('./infrastructure/logging');
const keyvault = require('./infrastructure/keyvault')
const config   = require('./config');

async function main() {
    try {
        logging.printConfiguration(config)
        
        logging.header(`Deploying Environment Instance Resources for ${config.env}`)

        //Environment resources
        await az.createResourceGroup(config.environmentInstanceResourceGroup)

        const azureContainerRepositoryPassword = await az.getAzureContainerRepositoryPassword(config.containerRepositoryName)
        await arm.deployEnvironmentInstanceResources(azureContainerRepositoryPassword)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
