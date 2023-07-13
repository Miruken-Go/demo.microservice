const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');
const git     = require('./git');

async function main() {
    try {
        console.log("Building teamsrv")
        config.requiredSecrets(['ghToken'])
        logging.printConfiguration(config)
        await az.login()

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
        await bash.execute(`
            az acr login -n ${config.containerRepositoryName}
        `)
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

        console.log("Built teamsrv")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
