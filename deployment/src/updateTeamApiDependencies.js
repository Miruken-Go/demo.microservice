const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredNonSecrets(['mirukenVersion'])
        config.requiredSecrets(['ghToken'])

        logging.printConfiguration(config)

        logging.header("Updating teamapi dependencies")

        await bash.execute(`
            cd teamapi
            go get github.com/miruken-go/miruken@${config.mirukenVersion}
        `)

        await git.commitAll(`Updated miruken to ${config.mirukenVersion}`)
        await git.push();

        // await bash.execute(`
        //     gh workflow run update-team-dependencies.yml \
        //         -f mirukenVersion=${config.mirukenVersion}      \
        //         -f teamapiVersion=${version}             \
        // `)

        console.log("Updated teamapi dependencies")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
