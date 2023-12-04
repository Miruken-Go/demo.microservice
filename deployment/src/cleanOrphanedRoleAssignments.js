import * as logging from '#infrastructure/logging.js'
import * as bash    from '#infrastructure/bash.js'
import { handle }   from '#infrastructure/handler.js'

handle(async () => {
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
})
