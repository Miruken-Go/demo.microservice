import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'mirukenVersion'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Updating team-api dependencies")

        await bash.execute(`
            cd team-api
            go get github.com/miruken-go/miruken@${variables.mirukenVersion}
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${variables.mirukenVersion}`)
            await git.push();

            await gh.sendRepositoryDispatch('updated-team-api', {})
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
