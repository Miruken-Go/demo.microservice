const az      = require('./az');
const arm     = require('./arm');
const logging = require('./logging');
const config  = require('./config');

async function main() {
    try {
        logging.printConfiguration(config)
        
        logging.header("Deploying B2C Configuration")

        


        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Script Failed")
    }
}

main()
