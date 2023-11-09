const logging            = require('./infrastructure/logging');
const b2c                = require('./infrastructure/b2c')
const b2cAppRegistration = require('./infrastructure/b2cAppRegistration')
const keyvault           = require('./infrastructure/keyvault')
const config             = require('./config');

async function main() {
    try {
        config.requiredEnvFileNonSecrets([
            'b2cDeploymentPipelineClientId',
            'authorizationServiceUsername',
        ])
        await keyvault.requireSecrets([
            'b2cDeploymentPipelineClientSecret',
        ])
        logging.printConfiguration(config)
        
        logging.header(`Deploying Manual Resource Configuration for ${config.env}`)

        await b2c.configureCustomPolicies()
        await b2cAppRegistration.configure()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
