function header (text) {
    const separator = "*********************************************************************"
    console.log(separator)
    console.log(text)
    console.log(separator)
}

function printConfiguration (config) {
    header('Configuration')

    const ignore = [
        'requiredSecrets',
        'requiredNonSecrets'
    ]

    for (const [key, value] of Object.entries(config)) {
        if (ignore.includes(key)){
            //ignore
        } else if (key == 'secrets') {
            console.log(`    ${key}:`);
            for (const [secretkey, secretvalue] of Object.entries(config.secrets)) {
                console.log(`        ${secretkey}: length ${secretvalue.length}`);
            }
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

module.exports = {
    header,
    printConfiguration 
}