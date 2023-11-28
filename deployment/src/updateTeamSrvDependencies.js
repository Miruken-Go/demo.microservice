const bash          = require('./infrastructure/bash')
const logging       = require('./infrastructure/logging');
const git           = require('./infrastructure/git');
const { variables } = require('./infrastructure/envVariables')

variables.requireEnvVariables([
    'mirukenVersion',
    'teamapiVersion',
    'teamVersion'
])

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Updating teams-rv dependencies")

        await bash.execute(`
            cd team-srv
            go get                                                                           \
                github.com/miruken-go/miruken@${variables.mirukenVersion}                    \
                github.com/miruken-go/demo.microservice/team-api@${variables.teamapiVersion} \
                github.com/miruken-go/demo.microservice/team@${variables.teamVersion}        \
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${variables.mirukenVersion}, teamapi to ${variables.teamapiVersion} and team to ${variables.teamVersion}`)
            await git.push();

            if (!variables.skipGitHubAction) {
                await bash.execute(`
                    gh workflow run build-team-srv.yml
                `)
            }
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
