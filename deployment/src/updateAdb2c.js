import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import * as gh       from '#infrastructure/gh.js'
import { handle }    from '#infrastructure/handler.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'mirukenVersion'
])

handle(async () => {
    logging.printEnvironmentVariables(variables)

    logging.header("Updating adb2c-api-connector-srv dependencies")

    await bash.execute(`
        cd adb2c
        go get github.com/miruken-go/miruken@${variables.mirukenVersion}
    `)

    if (await git.anyChanges()) {
        await git.commitAll(`Updated miruken to ${variables.mirukenVersion}`)
        await git.push();

        await gh.sendRepositoryDispatch('updated-adb2c', {})
    }
})
