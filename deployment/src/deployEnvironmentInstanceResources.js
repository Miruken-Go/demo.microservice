const az       = require('./az');
const arm      = require('./arm');
const logging  = require('./logging');
const config   = require('./config');
const b2c      = require('./b2c')
const keyvault = require('./keyvault')

async function main() {
    try {
        logging.printConfiguration(config)
        await keyvault.requireSecrets()
        
        logging.header("Deploying Environment")

        //Environment resources
        await az.createResourceGroup(config.resourceGroup)

        const getAzureContainerRepositoryPassword = await az.getAzureContainerRepositoryPassword(config.containerRepositoryName)
        await arm.deployEnvironmentInstanceResources(getAzureContainerRepositoryPassword)
        await b2c.configure()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
