const bash          = require('./infrastructure/bash')
const logging       = require('./infrastructure/logging');
const git           = require('./infrastructure/git');
const { variables } = require('./infrastructure/envVariables')

variables.requireEnvVariables([
    'mirukenVersion',
    'teamapiVersion'
])

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Updating team dependencies")

        await bash.execute(`
            cd team
            go get                                                                           \
                github.com/miruken-go/miruken@${variables.mirukenVersion}                    \
                github.com/miruken-go/demo.microservice/team-api@${variables.teamapiVersion} \
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${variables.mirukenVersion} and team-api to ${variables.teamapiVersion}`)
            await git.push();

            if (!variables.skipGitHubAction) {
                await bash.execute(`
                    gh workflow run build-team.yml
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
