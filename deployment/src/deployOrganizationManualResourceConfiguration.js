const logging          = require('./infrastructure/logging');
const { B2C }          = require('./infrastructure/b2c')
const { variables }    = require('./infrastructure/envVariables')
const { organization } = require('./config');

variables.requireEnvVariables([
    'env'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        logging.header(`Deploying Manual Resource Configuration for ${variables.env}`)

        const b2c = new B2C(organization)

        await b2c.configureCustomPolicies()
        await b2c.configureAppRegistrations()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
