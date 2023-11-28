const az               = require('./infrastructure/az');
const bash             = require('./infrastructure/bash')
const logging          = require('./infrastructure/logging');
const git              = require('./infrastructure/git');
const { organization } = require('./config');

async function tagContainerImageAndPush(imageName, app) {
    const appImage = `${app.imageName}:default`
    console.log(`imageName:    [${appImage}]`)
    await bash.execute(`
        docker tag ${imageName} ${appImage}
        docker push ${appImage}
    `)
}

async function main() {
    try {
        logging.printOrganization(organization)

        logging.header("Building defaultContainerImage")

        const version = `v${Math.floor(Date.now()/1000)}`.trim()
        const tag     = `defaultContainerImage/${version}`
        const imageName = 'defaultcontainerimage:latest' 

        console.log(`version: [${version}]`)
        console.log(`tag:     [${tag}]`)

        await bash.execute(`
            docker build              \
                -t ${imageName}             \
                defaultContainerImage \
        `)

        await az.loginToACR(organization.containerRepositoryName)

        //Push the default container for all the configured apps
        for (const app of organization.applications) {
            await tagContainerImageAndPush(imageName, app)
        }
        for(const domain of organization.domains) {
            for(const app of domain.applications) {
                await tagContainerImageAndPush(imageName, app)
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
