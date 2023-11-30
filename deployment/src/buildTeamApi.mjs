import * as bash     from '#infrastructure/bash.mjs'
import * as logging  from '#infrastructure/logging.mjs'
import * as git      from '#infrastructure/git.mjs'
import * as go       from '#infrastructure/go.mjs'
import { variables } from '#infrastructure/envVariables.mjs'

variables.requireEnvVariables([
    'repositoryPath'
])

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
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
      
        if (!variables.skipGitHubAction) {
            await bash.execute(`
                gh workflow run update-team-dependencies.yml \
                    -f mirukenVersion=${mirukenVersion}      \
                    -f teamapiVersion=${version}             \
            `)
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
