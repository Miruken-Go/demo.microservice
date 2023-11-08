const logging  = require('./infrastructure/logging');
const b2c      = require('./infrastructure/b2c')
const keyvault = require('./infrastructure/keyvault')
const config   = require('./config');

async function main() {
    try {
        config.requiredEnvFileNonSecrets([
            // 'b2cDeploymentPipelineClientId',
            // 'identityExperienceFrameworkClientId',
            // 'proxyIdentityExperienceFrameworkClientId',
            // 'b2cDomainName',
            // 'wellKnownOpenIdConfigurationUrl',
            // 'authorizationServiceUrl',
            // 'authorizationServiceUsername',
        ])
        config.requiredSecrets([
            'b2cDeploymentPipelineClientSecret',
        ])
        
        logging.printConfiguration(config)
        await keyvault.requireSecrets()
        
        logging.header(`Deploying Manual Resource Configuration for ${config.env}`)

        //App whatever automated configuration we can to the manual resource group
        await b2c.configure()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
