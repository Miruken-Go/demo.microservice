import {
    EnvVariables,
    EnvSecrets,
    logging,
    AZ,
    bash
} from 'ci.cd'

describe('organization global resources', () => {
    let az: AZ = undefined

    beforeAll(async () => {
        const variables = new EnvVariables()
            .required([
                'tenantId',
                'subscriptionId',
                'deploymentPipelineClientId',
                'deploymentPipelineClientSecret',
            ])
            .variables
        logging.printVariables(variables)

        const secrets = new EnvSecrets()
            .require([
                'deploymentPipelineClientSecret',
            ])
            .secrets
        logging.printSecrets(secrets)

        logging.header("Verifying Organization Global Resources")

        az = new AZ({
            tenantId:                       variables.tenantId,
            subscriptionId:                 variables.subscriptionId,
            deploymentPipelineClientId:     variables.deploymentPipelineClientId,
            deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
        })
        await az.login()
    })

    describe('resource group', () => {
        it('has expected name', async () => {
            const exists = await bash.execute(`
                az group exists --name 'majorleaguemiruken-global'
            `)
            expect(Boolean(exists)).toEqual(true)
        })
    })

    describe('container repository', () => {
        it('has expected name', async () => {
            const respositories:{ name: string }[] = await bash.json(`
                az acr list  
            `)
            const found = respositories.find(x => x.name === 'majorleaguemirukenglobal')
            expect(found).toBeDefined()
        })
    })
})
