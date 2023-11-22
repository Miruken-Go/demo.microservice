const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const git              = require('./infrastructure/git');
const { variables }    = require('./infrastructure/envVariables')
const { secrets }      = require('./infrastructure/envSecrets')
const { organization } = require('./config');

secrets.require([
   'ghToken' 
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printEnvironmentSecrets(secrets)
        logging.printOrganization(organization)

        logging.header("Building adb2c-auth-srv")

        const app = organization.getApplicationByName('adb2c-auth-srv')

        const version      = `v${Math.floor(Date.now()/1000)}`.trim()
        const imageName    = `${app.imageName}:${version}`
        const gitTag       = `${app.name}/${version}`
        const appSourceUrl = `${organization.repository}/releases/tag/${gitTag}`

        console.log(`version:      [${version}]`)
        console.log(`imageName:    [${imageName}]`)
        console.log(`gitTag:       [${gitTag}]`)
        console.log(`appSourceUrl: [${appSourceUrl}]`)

        await bash.execute(`
            docker build                                   \
                --build-arg app_source_url=${appSourceUrl} \
                --build-arg app_version=${version}         \
                -t ${imageName}                            \
                adb2c-auth-srv                             \
        `)

        await az.loginToACR()
        
        await bash.execute(`
            docker push ${imageName}
        `)

        await git.tagAndPush(gitTag)

        await bash.execute(`
            gh workflow run deploy-team-srv.yml \
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
