const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredNonSecrets(['mirukenVersion', 'teamapiVersion'])
        config.requiredSecrets(['ghToken'])

        logging.printConfiguration(config)

        logging.header("Updating team dependencies")

        await bash.execute(`
            cd teamapi
            go get github.com/miruken-go/miruken@${config.mirukenVersion} github.com/miruken-go/demo.microservice/teamapi@${config.teamapiVersion}
        `)

        await git.commitAll(`Updated miruken to ${config.mirukenVersion} and teamapi to ${config.teamapiVersion}`)
        await git.push();

        await bash.execute(`
            gh workflow run build-team.yml
        `)

        console.log("Updated team dependencies")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()