import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import * as gh       from '#infrastructure/gh.js'
import { handle }    from '#infrastructure/handler.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'mirukenVersion',
    'securityJwtVersion',
    'validatesPlayVersion',
    'configKoanfVersion'
])

handle(async () => {
    logging.printEnvironmentVariables(variables)

    logging.header("Updating miruken dependencies")

    // adb2c imports security/jwt, validates/play, and config/koanf directly
    // (split out of the miruken root module) - team-api does not.
    await bash.execute(`
        cd adb2c
        go get                                                                        \
            github.com/miruken-go/miruken@${variables.mirukenVersion}                 \
            github.com/miruken-go/miruken/security/jwt@${variables.securityJwtVersion}       \
            github.com/miruken-go/miruken/validates/play@${variables.validatesPlayVersion}   \
            github.com/miruken-go/miruken/config/koanf@${variables.configKoanfVersion}       \
    `)

    await bash.execute(`
        cd team-api
        go get github.com/miruken-go/miruken@${variables.mirukenVersion}
    `)

    if (await git.anyChanges()) {
        await git.commitAll(`Updated miruken to ${variables.mirukenVersion}`)
        await git.push();

        await gh.sendRepositoryDispatch('updated-miruken-dependencies')
    }
})
