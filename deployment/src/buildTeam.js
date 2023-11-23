const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const go      = require('./infrastructure/go');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableNonSecrets(['repositoryPath'])
        logging.printConfiguration(config)

        logging.header("Building team")

        await bash.execute(`
            cd team
            go test ./...
        `)

        const rawVersion = await bash.execute(`
            docker run --rm -v "${config.repositoryPath}:/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=team/v
        `)

        const version = `v${rawVersion}`
        const tag     = `team/${version}`

        console.log(`version: [${version}]`)
        console.log(`tag:     [${tag}]`)

        await git.tagAndPush(tag)

        const mirukenVersion = await go.getModuleVersion('team', 'github.com/miruken-go/miruken')
        const teamApiVersion = await go.getModuleVersion('team', 'github.com/miruken-go/demo.microservice/team-api')

        console.log(`mirukenVersion: [${mirukenVersion}]`)
        console.log(`teamApiVersion: [${teamApiVersion}]`)
      
        await bash.execute(`
            gh workflow run update-team-srv-dependencies.yml \
                -f mirukenVersion=${mirukenVersion}         \
                -f teamapiVersion=${teamApiVersion}         \
                -f teamVersion=${version}                   \
        `)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
