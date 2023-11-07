const az       = require('./az');
const arm      = require('./arm');
const logging  = require('./logging');
const config   = require('./config');
const keyvault = require('./keyvault')

async function main() {
    try {
        logging.printConfiguration(config)
        await keyvault.requireSecrets()
        
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
