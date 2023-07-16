const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredSecrets(['ghToken'])
        config.requiredNonSecrets(['mirukenVersion'])
        logging.printConfiguration(config)

        logging.header("Updating teamapi dependencies")

        await bash.execute(`
            cd teamapi
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
