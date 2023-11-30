import * as logging  from '#infrastructure/logging.js'
import { B2C }       from '#infrastructure/b2c.js'
import { variables } from '#infrastructure/envVariables.js'

import { 
    configDirectory,
    organization 
} from './config.js'

variables.requireEnvVariables([
    'env'
])

variables.requireEnvFileVariables(configDirectory, [
    'b2cDeploymentPipelineClientId'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printOrganization(organization)

        logging.header(`Deploying Manual Resource Configuration for ${variables.env}`)

        const b2c = new B2C(organization, variables.b2cDeploymentPipelineClientId)

        await b2c.configureAppRegistrations()
        await b2c.configureCustomPolicies()

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
