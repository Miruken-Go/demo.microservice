import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import * as go       from '#infrastructure/go.js'
import * as gh       from '#infrastructure/gh.js'
import { handle }    from '#infrastructure/handler.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'repositoryPath'
])

handle(async () => {
    logging.printEnvironmentVariables(variables)

    logging.header("Building team-api")

    await bash.execute(`
        cd team-api
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

    await git.tagAndPush(tag)

    const mirukenVersion = await go.getModuleVersion('team-api', 'github.com/miruken-go/miruken')
    console.log(`mirukenVersion: [${mirukenVersion}]`)
    
    await gh.sendRepositoryDispatch('built-team-api', {
        mirukenVersion: mirukenVersion,
        teamapiVersion: version,
    })
})
