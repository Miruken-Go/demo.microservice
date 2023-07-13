const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredNonSecrets(['mirukenVersion', 'teamapiVersion', 'teamVersion'])
        config.requiredSecrets(['ghToken'])

        logging.printConfiguration(config)

        logging.header("Updating teamsrv dependencies")

        await bash.execute(`
            cd teamapi
            go get github.com/miruken-go/miruken@${config.mirukenVersion} github.com/miruken-go/demo.microservice/teamapi@${config.teamapiVersion} 	github.com/miruken-go/demo.microservice/team v0.2.11@${config.teamVersion}
        `)

        await git.commitAll(`Updated miruken to ${config.mirukenVersion}, teamapi to ${config.teamapiVersion} and team to ${config.teamVersion}`)
        await git.push();

        await bash.execute(`
            gh workflow run build-teamsrv.yml
        `)

        console.log("Updated teamsrv dependencies")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
