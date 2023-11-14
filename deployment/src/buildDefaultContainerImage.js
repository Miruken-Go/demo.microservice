const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const git              = require('./infrastructure/git');
const { organization } = require('./config');

async function main() {
    try {
        logging.printOrganization(organization)

        logging.header("Building defaultContainerImage")

        const version   = `v${Math.floor(Date.now()/1000)}`.trim()
        const imageName = `${organization.apiConnectorImageName}:default`
        const tag       = `defaultContainerImage/${version}`

        console.log(`version:      [${version}]`)
        console.log(`imageName:    [${imageName}]`)
        console.log(`tag:          [${tag}]`)

        await bash.execute(`
            docker build                                   \
                -t ${imageName}                            \
                defaultContainerImage                      \
        `)

        await az.loginToACR(organization.containerRepositoryName)

        await bash.execute(`
            docker push ${imageName}
        `)

        //Push a default container for all the apps
        for(const domain of organization.domains) {
            for (const app of domain.applications) {
                const appImage = `${app.imageName}:default`
                console.log(`imageName:    [${appImage}]`)
                await bash.execute(`
                    docker tag ${imageName} ${appImage}
                    docker push ${appImage}
                `)
            }
        }

        await git.tagAndPush(tag)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
