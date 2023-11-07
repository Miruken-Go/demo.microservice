const az      = require('./infrastructure/az');
const arm     = require('./infrastructure/arm');
const logging = require('./infrastructure/logging');
const config  = require('./config');

async function main() {
    try {
        logging.printConfiguration(config)

        logging.header("Deploying Global Resources")

        //Provider Registrations
        await az.registerAzureProvider('Microsoft.AzureActiveDirectory')
        await az.registerAzureProvider('Microsoft.App')
        await az.registerAzureProvider('Microsoft.OperationalInsights')

        //Global resources 
        await az.createResourceGroup(config.globalResourceGroup)
        await arm.deployGlobalResources()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
