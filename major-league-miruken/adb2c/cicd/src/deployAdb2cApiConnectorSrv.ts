import { config } from 'dotenv'
config({ path: `.env.${process.env.env}` })

import { 
    handle,
    logging,
    bash,
    EnvSecrets,
    EnvVariables,
    AZ,
} from 'ci.cd'

import { organization } from './domains'

handle(async () => {
    const variables = new EnvVariables()
        .required([
            'tag',
            'tenantId',
            'subscriptionId',
            'deploymentPipelineClientId',
            'authorizationServiceUsername'
        ])
        .variables

    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
        ])
        .secrets
    logging.printSecrets(secrets)

    const application = organization.getApplicationByName("adb2c-api-connector-srv")

    logging.header(`Deploying ${application.name}`)

    const envVars = [
        `Databases__Azure__ConnectionUri='secretref:cosmos-connection-string'`,
        `Login__Adb2c__0__Module='login.pwd'`,
        `Login__Adb2c__0__Options__Credentials__0__Username='${variables.authorizationServiceUsername}'`,
        `Login__Adb2c__0__Options__Credentials__0__Password='secretref:authorization-service-password'`,
    ]

    await new AZ({
        tenantId:                       variables.tenantId,
        subscriptionId:                 variables.subscriptionId,
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
    }).login()

    //https://learn.microsoft.com/en-us/cli/azure/containerapp?view=azure-cli-latest#az-containerapp-update
    //Create the new revision
    const now            = `${Math.floor(Date.now()/1000)}`.trim()
    const revisionSuffix = `${variables.tag}-${now}`
    await bash.execute(`
        az containerapp update                                \
            -n ${application.containerAppName}                \
            -g ${application.resourceGroups.instance}         \
            --image ${application.imageName}:${variables.tag} \
            --container-name ${application.name}              \
            --revision-suffix ${revisionSuffix}               \
            --replace-env-vars ${envVars.join(' ')}           \
    `)

    const revisions = await bash.json(`
        az containerapp revision list                 \
            -n ${application.containerAppName}        \
            -g ${application.resourceGroups.instance} \
    `)

    const revisionsToActivate   = revisions.filter(r => r.name.includes(revisionSuffix))
    const revisionsToDeactivate = revisions.filter(r => !r.name.includes(revisionSuffix))

    //You must have an active revision before deactivating the rest
    for (const revision of revisionsToActivate) {
        if (revision.properties.active !== true) {
            console.log(`Activating: ${revision.name}`)
            await bash.execute(`
                az containerapp revision activate             \
                    -n ${application.containerAppName}        \
                    -g ${application.resourceGroups.instance} \
                    --revision ${revision.name}               \
            `)
        }

        await bash.execute(`
            az containerapp ingress traffic set           \
                -n ${application.containerAppName}        \
                -g ${application.resourceGroups.instance} \
                --revision-weight ${revision.name}=100    \
        `)
    }

    for (const revision of revisionsToDeactivate) {
        if (revision.properties.active === true) {
            console.log(`Dectivate: ${revision.name}`)
            await bash.execute(`
                az containerapp revision deactivate           \
                    -n ${application.containerAppName}        \
                    -g ${application.resourceGroups.instance} \
                    --revision ${revision.name}               \
            `)
        }
    }
})
