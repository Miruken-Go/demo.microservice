const az      = require('./infrastructure/az');
const bash    = require('./infrastructure/bash')
const logging = require('./infrastructure/logging');
const git     = require('./infrastructure/git');
const config  = require('./config');

async function main() {
    try {
        config.requiredEnvironmentVariableSecrets(['ghToken'])
        logging.printConfiguration(config)

        logging.header("Building teamsrv")

        const version      = `v${Math.floor(Date.now()/1000)}`.trim()
        const imageName    = `${config.imageName}:${version}`
        const tag          = `${config.appName}/${version}`
        const appSourceUrl = `${config.repository}/releases/tag/${tag}`

        console.log(`version:      [${version}]`)
        console.log(`imageName:    [${imageName}]`)
        console.log(`tag:          [${tag}]`)
        console.log(`appSourceUrl: [${appSourceUrl}]`)

        await bash.execute(`
            docker build                                   \
                --build-arg app_source_url=${appSourceUrl} \
                --build-arg app_version=${version}         \
                -t ${imageName}                            \
                teamsrv                                    \
        `)

        await az.loginToACR()
        
        await bash.execute(`
            docker push ${imageName}
        `)

        await git.tagAndPush(tag)

        await bash.execute(`
            gh workflow run deploy-teamsrv.yml \
                -f env=dev                     \
                -f instance=ci                 \
                -f tag=${version}              \
        `)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
