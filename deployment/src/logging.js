function header (text) {
    const separator = "********************************************"
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
        } else if (key == 'armParams') {
            console.log(`    ${key}:`);
            for (const [armkey, armvalue] of Object.entries(config.armParams)) {
                console.log(`        ${armkey}: ${armvalue.value}`);
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