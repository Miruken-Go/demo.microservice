const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableSecrets(['ghToken'])
        config.requiredEnvironmentVariableNonSecrets(['mirukenVersion'])
        logging.printConfiguration(config)

        logging.header("Updating team-api dependencies")

        await bash.execute(`
            cd team-api
            go get github.com/miruken-go/miruken@${config.mirukenVersion}
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${config.mirukenVersion}`)
            await git.push();

            await bash.execute(`
                gh workflow run build-teamapi.yml
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
