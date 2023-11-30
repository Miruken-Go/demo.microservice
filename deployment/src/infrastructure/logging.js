import { inspect } from 'node:util'

export function header (text) {
    const separator = "*********************************************************************"
    console.log(separator)
    console.log(text)
    console.log(separator)
}

export function printObject (text, object) {
    header(text)
    console.log(inspect(object, { depth: null }))
}

export function printOrganization (object) {
    printObject("Organization Configuration", object)
}

export function printEnvironmentVariables (config) {
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

export function printEnvironmentSecrets(config) {
    header('Environment Secrets')

    for (const [key, value] of Object.entries(config)) {
        if (typeof value === 'function') {
            //ignore
        } else {
            console.log(`    ${key}: length ${value.length}`);
        }
    }
}