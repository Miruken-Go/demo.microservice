import * as logging from '#infrastructure/logging.js'
import * as bash    from '#infrastructure/bash.js'

async function main() {
    try {
        logging.header("Cleaning up orphaned application security principals")

        const ids = await bash.json(`
            az role assignment list --all --query "[?principalName==''].id"    
        `)

        if (ids.length) {
            //console.log(`az role assignment delete --ids "${ids.join(' ')}"`)
            // await bash.execute(`
            //     az role assignment delete --ids "${ids.join(' ')}"
            // `)

            // console.log(`Deleted ${ids.length} orphaned application security principals`)
        }

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
