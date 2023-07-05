const az      = require('./az');
const arm     = require('./arm');
const logging = require('./logging');
const config  = require('./config');

async function main() {
    try {
        console.log("Deploying Shared Environment")
        logging.printConfiguration(config)
        await az.login()

        //Provider Registrations
        await az.registerAzureProvider('Microsoft.AzureActiveDirectory')
        await az.registerAzureProvider('Microsoft.App')
        await az.registerAzureProvider('Microsoft.OperationalInsights')

        //Shared resources 
        await az.createResourceGroup(config.sharedResourceGroup)
        await arm.deploySharedResources()

        console.log("Configuration Succeded")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Configuration Failed")
    }
}

main()
