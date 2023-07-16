const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');
const go      = require('./go');

async function main() {
    try {
        config.requiredSecrets(['ghToken'])
        config.requiredNonSecrets(['repositoryPath'])
        logging.printConfiguration(config)

        logging.header("Building teamapi")

        await bash.execute(`
            cd teamapi
            go test ./...
        `)

        const rawVersion = await bash.execute(`
            docker run --rm -v "${config.repositoryPath}:/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=teamapi/v
        `)

        const version = `v${rawVersion}`
        const tag     = `teamapi/${version}`

        console.log(`version: [${version}]`)
        console.log(`tag:     [${tag}]`)

        await git.tagAndPush(tag)

        const mirukenVersion = await go.getModuleVersion('teamapi', 'github.com/miruken-go/miruken')
        console.log(`mirukenVersion: [${mirukenVersion}]`)
      
        await bash.execute(`
            gh workflow run update-team-dependencies.yml \
                -f mirukenVersion=${mirukenVersion}      \
                -f teamapiVersion=${version}             \
        `)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
