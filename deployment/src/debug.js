const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredSecrets(['ghToken'])
        logging.printConfiguration(config)

        logging.header("Debugging")

        const rawVersion = await bash.execute(`
            docker run --rm -v "$(pwd):/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=teamapi/v
        `)

        console.log(`rawVersion: ${rawVersion}`)

        console.log("Debugging Complete")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
