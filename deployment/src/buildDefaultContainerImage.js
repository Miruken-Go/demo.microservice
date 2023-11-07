const az      = require('./infrastructure/az');
const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableSecrets(['ghToken'])
        logging.printConfiguration(config)

        logging.header("Building defaultContainerImage")

        const version      = `v${Math.floor(Date.now()/1000)}`.trim()
        const imageName    = `${config.imageName}:default`
        const tag          = `${config.defaultContainerImage}/${version}`

        console.log(`version:      [${version}]`)
        console.log(`imageName:    [${imageName}]`)
        console.log(`tag:          [${tag}]`)

        await bash.execute(`
            docker build                                   \
                -t ${imageName}                            \
                defaultContainerImage                      \
        `)

        await az.loginToACR()

        await bash.execute(`
            docker push ${imageName}
        `)

        await git.tagAndPush(tag)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
