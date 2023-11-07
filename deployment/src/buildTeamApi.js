const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const go      = require('./infrastructure/go');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableSecrets(['ghToken'])
        config.requiredEnvironmentVariableNonSecrets(['repositoryPath'])
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
