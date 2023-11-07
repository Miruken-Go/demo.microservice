const logging  = require('./logging');
const config   = require('./config');
const b2c      = require('./b2c')
const keyvault = require('./keyvault')

async function main() {
    try {
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
