import * as az          from '#infrastructure/az.js'
import * as bash        from '#infrastructure/bash.js'
import * as logging     from '#infrastructure/logging.js'
import * as git         from '#infrastructure/git.js'
import { handle }       from '#infrastructure/handler.js'
import { organization } from './config.js'

async function tagContainerImageAndPush(imageName, app) {
    const appImage = `${app.imageName}:default`
    console.log(`imageName:    [${appImage}]`)
    await bash.execute(`
        docker tag ${imageName} ${appImage}
        docker push ${appImage}
    `)
}

handle(async () => {
    logging.printDomain(organization)

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

    await az.loginToACR(organization.containerRepository.name)

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
})
