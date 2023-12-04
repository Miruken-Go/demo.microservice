import * as bash     from '#infrastructure/bash.js'
import * as logging  from '#infrastructure/logging.js'
import * as git      from '#infrastructure/git.js'
import { variables } from '#infrastructure/envVariables.js'

variables.requireEnvVariables([
    'mirukenVersion',
    'teamapiVersion',
    'teamVersion'
])

async function main() {
    try {
        logging.printEnvironmentVariables(variables)

        logging.header("Updating teams-rv dependencies")

        await bash.execute(`
            cd team-srv
            go get                                                                           \
                github.com/miruken-go/miruken@${variables.mirukenVersion}                    \
                github.com/miruken-go/demo.microservice/team-api@${variables.teamapiVersion} \
                github.com/miruken-go/demo.microservice/team@${variables.teamVersion}        \
        `)

        if (await git.anyChanges()) {
            await git.commitAll(`Updated miruken to ${variables.mirukenVersion}, teamapi to ${variables.teamapiVersion} and team to ${variables.teamVersion}`)
            await git.push();

            await gh.sendRepositoryDispatch('updated-team-srv-dependencies', {})
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
