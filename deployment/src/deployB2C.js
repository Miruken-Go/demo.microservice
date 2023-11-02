const az      = require('./az');
const arm     = require('./arm');
const logging = require('./logging');
const config  = require('./config');

const graph   = require('./graph');
const fs      = require('fs')
const path     = require('path')

async function main() {

    try {
        logging.printConfiguration(config)
        
        logging.header("Deploying B2C Configuration")

        console.log('Updating Custom Policies')

        const customPoliciesDirectory = path.resolve(__dirname, '../custom-policies')
        const customPoliciesFileOrder = [
            'TrustFrameworkBase.xml',
            'TrustFrameworkLocalization.xml',
            'TrustFrameworkExtensions.xml',
            'SignUpOrSignin.xml',
            'ProfileEdit.xml',
            'PasswordReset.xml',
        ]
        for (const file of customPoliciesFileOrder) {
                const policyId = `B2C_1A_${path.basename(file, '.xml')}`
                const filePath = path.join(customPoliciesDirectory, file)
                let xml        = fs.readFileSync(filePath,{encoding: 'utf-8'}) 

                xml = xml.replace(/{B2C_DOMAIN_NAME}/g,                              config.b2cDomainName)
                xml = xml.replace(/{IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g,       config.identityExperienceFrameworkClientId)
                xml = xml.replace(/{PROXY_IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g, config.proxyIdentityExperienceFrameworkClientId)

                await graph.updateTrustFrameworkPolicy(policyId, xml)
        };

        console.log("Script completed successfully")
    } catch (error) {
        process.exitCode = 1
        //console.log(error)
        console.log("Script Failed")
    }
}

main()
