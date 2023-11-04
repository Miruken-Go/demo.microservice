const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        config.requiredSecrets(['ghToken'])
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
                --build-arg app_version=${version}         \
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
