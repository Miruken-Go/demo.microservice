const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredSecrets(['ghToken'])
        logging.printConfiguration(config)

        logging.header("Building teamapi")

        await bash.execute(`
            docker run --rm -v $(pwd):/go/src --workdir=/go/src/teamapi golang:1.20 go test ./...
        `)

        const rawVersion = await bash.execute(`
            docker run --rm -v "$(pwd):/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=teamapi/v
        `)

        const version = `v${rawVersion}`
        const tag     = `teamapi/${version}`

        console.log(`version: [${version}]`)
        console.log(`tag:     [${tag}]`)

        await git.tagAndPush(tag)

        const mirukenVersion = await bash.execute(`
            docker run                                                                 \
                -v $(pwd):/go/src                                                      \
                --workdir=/go/src/teamapi                                              \
                golang:1.20                                                            \
                go list -m all | grep github.com/miruken-go/miruken | awk '{print $2}' \
        `)
        console.log(`mirukenVersion: [${mirukenVersion}]`)
      
        await bash.execute(`
            gh workflow run update-team-dependencies.yml \
                -f mirukenVersion=${mirukenVersion}      \
                -f teamapiVersion=${version}             \
        `)

        console.log("Built teamapi")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Deployment Failed")
    }
}

main()