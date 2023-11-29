const bash          = require('./infrastructure/bash')
const logging       = require('./infrastructure/logging');
const git           = require('./infrastructure/git');
const go            = require('./infrastructure/go');
const { variables } = require('./infrastructure/envVariables')

variables.requireEnvVariables([
    'repositoryPath'
])

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Building team")

        await bash.execute(`
            cd team
            go test ./...
        `)

        //This docker container is running docker in docker from github actions
        //Therefore using $(pwd) to get the working directory would be the working directory of the running container 
        //Not the working directory from the host system. So we need to pass in the repository path.
        const rawVersion = await bash.execute(`
            docker run --rm -v "${variables.repositoryPath}:/repo" \
            gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=team/v
        `)

        const version = `v${rawVersion}`
        const tag     = `team/${version}`

        console.log(`version: [${version}]`)
        console.log(`tag:     [${tag}]`)

        await git.tagAndPush(tag)

        const mirukenVersion = await go.getModuleVersion('team', 'github.com/miruken-go/miruken')
        const teamApiVersion = await go.getModuleVersion('team', 'github.com/miruken-go/demo.microservice/team-api')

        console.log(`mirukenVersion: [${mirukenVersion}]`)
        console.log(`teamApiVersion: [${teamApiVersion}]`)
      
        if (!variables.skipGitHubAction) {
            await bash.execute(`
                gh workflow run update-team-srv-dependencies.yml \
                    -f mirukenVersion=${mirukenVersion}         \
                    -f teamapiVersion=${teamApiVersion}         \
                    -f teamVersion=${version}                   \
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
