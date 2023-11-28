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

        logging.header("Building team-api")

        await bash.execute(`
            cd team-api
            go test ./...
        `)

        const rawVersion = await bash.execute(`
            docker run                                     \
                --rm                                       \
                -v "${variables.repositoryPath}:/repo"     \
                gittools/gitversion:5.12.0-alpine.3.14-6.0 \
                    /repo                                  \
                    /showvariable SemVer                   \
                    /overrideconfig tag-prefix=team-api/v  \
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
