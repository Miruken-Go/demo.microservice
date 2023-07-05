const util = require('util')
const exec = util.promisify(require('child_process').exec);

async function execute(cmd, suppressLog) { 
    const { stdout, stderr } = await exec(cmd);
    if (!suppressLog) {
        console.log('bash stdout:', stdout);
        if (stderr){
            console.log('bash stderr:', stderr);
        }
    }
    return stdout;
}

async function json(cmd, suppressLog) { 
    const response = await execute(cmd, suppressLog)
    return JSON.parse(response)
}

module.exports = {
    execute,
    json
}
