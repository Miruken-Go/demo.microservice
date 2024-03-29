import * as logging  from '#infrastructure/logging.js'
import * as gh       from '#infrastructure/gh.js'
import { handle }    from '#infrastructure/handler.js'
import { B2C }       from '#infrastructure/b2c.js'
import { Users }     from '#infrastructure/users.js'
import { variables } from '#infrastructure/envVariables.js'
import path          from 'node:path'

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

handle(async () => {
    logging.printEnvironmentVariables(variables)
    logging.printDomain(organization)

    logging.header(`Deploying Manual Resource Configuration for ${variables.env}`)

    const b2c = new B2C(organization, variables.b2cDeploymentPipelineClientId)

    await b2c.configureAppRegistrations()
    await b2c.configureCustomPolicies(path.join(configDirectory, 'custom-policies'))

    const users = new Users(organization, variables.b2cDeploymentPipelineClientId)
    await users.configureBootstrapUsers()

    await gh.sendRepositoryDispatch(`provisioned-organization-environment-manual-resource-configuration`, {
        env:      organization.env,
        instance: organization.instance,
    })
})
