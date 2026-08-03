import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import * as gh       from '#infrastructure/gh.js'
import { handle }    from '#infrastructure/handler.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'mirukenVersion',
    'teamapiVersion',
    'securityJwtVersion',
    'validatesPlayVersion',
    'configKoanfVersion'
])

handle(async () => {
    logging.printEnvironmentVariables(variables)

    logging.header("Updating team dependencies")

    // team imports security/jwt, validates/play, and config/koanf directly
    // (split out of the miruken root module).
    await bash.execute(`
        cd team
        go get                                                                                \
            github.com/miruken-go/miruken@${variables.mirukenVersion}                          \
            github.com/miruken-go/miruken/security/jwt@${variables.securityJwtVersion}         \
            github.com/miruken-go/miruken/validates/play@${variables.validatesPlayVersion}     \
            github.com/miruken-go/miruken/config/koanf@${variables.configKoanfVersion}         \
            github.com/miruken-go/demo.microservice/team-api@${variables.teamapiVersion}       \
    `)

    if (await git.anyChanges()) {
        await git.commitAll(`Updated miruken to ${variables.mirukenVersion} and team-api to ${variables.teamapiVersion}`)
        await git.push();

        await gh.sendRepositoryDispatch('updated-team-dependencies')
    }
})
