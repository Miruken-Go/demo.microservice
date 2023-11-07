const az       = require('./infrastructure/az');
const arm      = require('./infrastructure/arm');
const logging  = require('./infrastructure/logging');
const config   = require('./config');

async function main() {
    try {
        logging.printConfiguration(config)
        
        logging.header("Deploying Common Environment Resources")

        await az.createResourceGroup(config.commonEnvironmentResourceGroup)
        await az.createResourceGroup(config.manualEnvironmentResourceGroup)
        await arm.deployCommonEnvironmentResources()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
