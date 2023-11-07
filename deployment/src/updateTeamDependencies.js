const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableNonSecrets(['mirukenVersion', 'teamapiVersion'])
        config.requiredEnvironmentVariableSecrets(['ghToken'])
        logging.printConfiguration(config)

        logging.header("Updating team dependencies")

        await bash.execute(`
            cd team
            go get github.com/miruken-go/miruken@${config.mirukenVersion} github.com/miruken-go/demo.microservice/teamapi@${config.teamapiVersion}
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${config.mirukenVersion} and teamapi to ${config.teamapiVersion}`)
            await git.push();

            await bash.execute(`
                gh workflow run build-team.yml
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
