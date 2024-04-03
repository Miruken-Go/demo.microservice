import { 
    handle,
    logging,
    bash,
    EnvSecrets,
    EnvVariables,
    go,
    Git,
    GH
} from 'ci.cd'

handle(async() => {
    const variables = new EnvVariables()
        .required([
            'repositoryPath',
            'repository',
            'repositoryOwner',
            'ref',
        ])
        .optional(['skipRepositoryDispatches'])
        .variables
    logging.printVariables(variables)

    const secrets = new EnvSecrets()
        .require(['GH_TOKEN'])
        .secrets
    logging.printSecrets(secrets)

    logging.header("Publishing team-api")

    await bash.execute(`
        cd ../app
        go test ./...
    `)

    //This docker container is running docker in docker from github actions
    //Therefore using $(pwd) to get the working directory would be the working directory of the running container 
    //Not the working directory from the host system. So we need to pass in the repository path.
    const rawVersion = await bash.execute(`
        docker run --rm -v '${variables.repositoryPath}:/repo' \
        gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=team-api/v
    `)

    const version = `v${rawVersion}`
    const tag     = `team-api/${version}`

    console.log(`version: [${version}]`)
    console.log(`tag:     [${tag}]`)

    await new Git(secrets.GH_TOKEN)
        .tagAndPush(tag)

    const mirukenVersion = await go.getModuleVersion('team-api', 'github.com/miruken-go/miruken')
    console.log(`mirukenVersion: [${mirukenVersion}]`)
    
    await new GH({
        ghToken:                  secrets.GH_TOKEN,
        ref:                      variables.ref,
        repository:               variables.repository,
        repositoryOwner:          variables.repositoryOwner,
        skipRepositoryDispatches: Boolean(variables.skipRepositoryDispatches)
    }).sendRepositoryDispatch('built-team-api', {
        mirukenVersion: mirukenVersion,
        teamapiVersion: version,
    }, variables.repository)
})