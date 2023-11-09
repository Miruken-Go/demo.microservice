const config             = require('./config');
const az                 = require('./infrastructure/az');
const keyvault           = require('./infrastructure/keyvault')
const b2cAppRegistration = require('./infrastructure/b2cAppRegistration') 
const {ApplicationType}  = require('./infrastructure/systemDescription')

async function main() {
    try {
        const app = await az.getContainerAppUrl('teamsrv-dev-ci')
        console.log(app)

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
