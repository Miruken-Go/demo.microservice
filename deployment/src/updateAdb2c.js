const bash          = require('./infrastructure/bash')
const logging       = require('./infrastructure/logging');
const git           = require('./infrastructure/git');
const { variables } = require('./infrastructure/envVariables')

variables.requireEnvVariables([
    'mirukenVersion'
])

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Updating adb2c-api-connector-srv dependencies")

        await bash.execute(`
            cd adb2c
            go get github.com/miruken-go/miruken@${variables.mirukenVersion}
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${variables.mirukenVersion}`)
            await git.push();

            if (!variables.skipGitHubAction) {
                await bash.execute(`
                    gh workflow run build-adb2c-api-connector-srv.yml
                    gh workflow run build-adb2c-auth-srv.yml
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
