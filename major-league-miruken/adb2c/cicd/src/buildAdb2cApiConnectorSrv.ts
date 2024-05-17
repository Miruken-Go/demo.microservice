import { 
    handle,
    logging,
    bash,
    EnvSecrets,
    EnvVariables,
    AZ,
    Git,
    GH
} from 'ci.cd'

import { organization } from './domains'

handle(async () => {
    const variables = new EnvVariables()
        .required([
            'tenantId',
            'subscriptionId',
            'deploymentPipelineClientId',
            'deploymentPipelineClientSecret',
            'repository',
            'repositoryOwner',
            'ref',
        ])
        .optional(['skipRepositoryDispatches'])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require([
            'deploymentPipelineClientSecret',
            'GH_TOKEN'
        ])
        .secrets
    logging.printSecrets(secrets)

    const appName = 'adb2c-api-connector-srv'

    const app = organization.getApplicationByName(appName)

    const version      = `v${Math.floor(Date.now()/1000)}`.trim()
    const imageName    = `${app.imageName}:${version}`
    const gitTag       = `${app.name}/${version}`
    const appSourceUrl = `${organization.gitRepositoryUrl}/releases/tag/${gitTag}`

    console.log(`version:      [${version}]`)
    console.log(`imageName:    [${imageName}]`)
    console.log(`gitTag:       [${gitTag}]`)
    console.log(`appSourceUrl: [${appSourceUrl}]`)

    await bash.execute(`
        cd ../../
        docker build                                   \
            --progress plain                           \
            --build-arg app_source_url=${appSourceUrl} \
            --build-arg app_version=${version}         \
            -f adb2c/cmd/api-connector-srv/Dockerfile  \
            -t ${imageName}                            \
            .                                          \
    `)

    await new AZ({
        tenantId:                       variables.tenantId,
        subscriptionId:                 variables.subscriptionId,
        deploymentPipelineClientId:     variables.deploymentPipelineClientId,
        deploymentPipelineClientSecret: secrets.deploymentPipelineClientSecret
    }).loginToACR(organization.resources.containerRepository.name)
    
    await bash.execute(`
        docker push ${imageName}
    `)

    await new Git(secrets.GH_TOKEN)
        .tagAndPush(gitTag)

    await new GH({
        ghToken:                  secrets.GH_TOKEN,
        ref:                      variables.ref,
        repository:               variables.repository,
        repositoryOwner:          variables.repositoryOwner,
        skipRepositoryDispatches: Boolean(variables.skipRepositoryDispatches)
    }).sendRepositoryDispatch('built-team-srv', {
        env:      'dev',
        instance: 'ci',
        tag:      version,
    }, variables.repository)
})
