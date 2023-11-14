const { inspect } = require('node:util');

function header (text) {
    const separator = "*********************************************************************"
    console.log(separator)
    console.log(text)
    console.log(separator)
}

function printOrganization (object) {
    header("Organization Configuration")
    console.log(inspect(object, { depth: null }))
}

function printEnvironmentVariables (config) {
    header('Environment Variables')

    for (const [key, value] of Object.entries(config)) {
        if (typeof value === 'function') {
            //ignore
        } else if (Array.isArray(value)) {
            console.log(`    ${key}:`);
            for (const [_, arrayValue] of Object.entries(value)) {
                console.log(`        ${arrayValue}`);
            }
        } else {
            console.log(`    ${key}: ${value}`);
        }
    }
}

function printEnvironmentSecrets(config) {
    header('Environment Secrets')

    for (const [key, value] of Object.entries(config)) {
        if (typeof value === 'function') {
            //ignore
        } else {
            console.log(`    ${key}: length ${value.length}`);
        }
    }
}

module.exports = {
    header,
    printOrganization,
    printEnvironmentVariables,
    printEnvironmentSecrets,
}