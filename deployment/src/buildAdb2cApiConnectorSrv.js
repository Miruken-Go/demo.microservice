import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as logging     from '#infrastructure/logging.js'
import * as git         from '#infrastructure/git.js'
import { variables }    from '#infrastructure/envVariables.js'
import { secrets }      from '#infrastructure/envSecrets.js'
import { organization } from './config.js'

variables.optionalEnvVariables([
    'skipGitHubAction'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)
        logging.printEnvironmentSecrets(secrets)
        logging.printOrganization(organization)

        const appName = 'adb2c-api-connector-srv'

        logging.header(`Building ${appName}`)

        const app = organization.getApplicationByName(appName)

        const version      = `v${Math.floor(Date.now()/1000)}`.trim()
        const imageName    = `${app.imageName}:${version}`
        const gitTag       = `${app.name}/${version}`
        const appSourceUrl = `${organization.gitRepositoryUrl}/releases/tag/${gitTag}`

        console.log(`version:      [${version}]`)
        console.log(`imageName:    [${imageName}]`)
        console.log(`gitTag:       [${gitTag}]`)
        console.log(`appSourceUrl: [${appSourceUrl}]`)

        await bash.execute(`
            docker build                                   \
                --progress plain                           \
                --build-arg app_source_url=${appSourceUrl} \
                --build-arg app_version=${version}         \
                -f adb2c/cmd/api-connector-srv/Dockerfile  \
                -t ${imageName}                            \
                .                                          \
        `)

        await az.loginToACR(organization.containerRepositoryName)
        
        await bash.execute(`
            docker push ${imageName}
        `)

        await git.tagAndPush(gitTag)

        if (!variables.skipGitHubAction) {
            await bash.execute(`
                gh workflow run deploy-${appName}.yml \
                    -f env=dev                     \
                    -f instance=ci                 \
                    -f tag=${version}              \
            `)
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
