const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');

async function main() {
    try {
        console.log("Deploying teamsrv")
        config.requiredNonSecrets([
            'tag'
        ])
        logging.printConfiguration(config)
        await az.login()

        //Create the new revision
        await bash.execute(`
            az containerapp update                        \
                -n ${config.prefix}                       \
                -g ${config.resourceGroup}                \
                --image ${config.imageName}:${config.tag} \
                --container-name ${config.appName}        \
                --revision-suffix ${config.tag}           \
        `)

        const revisions = await bash.json(`
            az containerapp revision list  \
                -n ${config.prefix}        \
                -g ${config.resourceGroup} \
        `)

        const revisionsToActivate   = revisions.filter(r => r.name.includes(config.tag))
        const revisionsToDeactivate = revisions.filter(r => !r.name.includes(config.tag))

        //You must have an active revision before deactivating the rest
        for (const revision of revisionsToActivate) {
            if (revision.properties.active !== true) {
                console.log(`Activating: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision activate \
                        -n ${config.prefix}           \
                        -g ${config.resourceGroup}    \
                        --revision ${revision.name}   \
                `)
            }

            await bash.execute(`
                az containerapp ingress traffic set \
                    -n ${config.prefix}           \
                    -g ${config.resourceGroup}    \
                    --revision-weight ${revision.name}=100
            `)
        }

        for (const revision of revisionsToDeactivate) {
            if (revision.properties.active === true) {
                console.log(`Dectivate: ${revision.name}`)
                await bash.execute(`
                    az containerapp revision deactivate \
                        -n ${config.prefix}           \
                        -g ${config.resourceGroup}    \
                        --revision ${revision.name}   \
                `)
            }
        }

        console.log("Deployment Succeded")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Deployment Failed")
    }
}

main()
