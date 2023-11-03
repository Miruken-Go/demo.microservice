const logging = require('./logging');
const config  = require('./config');
const graph   = require('./graph');
const fs      = require('fs')
const path    = require('path')

async function configure () {

    logging.header("Deploying B2C Configuration")

    //https://learn.microsoft.com/en-us/azure/active-directory-b2c/deploy-custom-policies-devops
    const customPoliciesDirectory = path.resolve(__dirname, '../custom-policies')
    const customPoliciesFileOrder = [
        'TrustFrameworkBase.xml',
        'TrustFrameworkLocalization.xml',
        'TrustFrameworkExtensions.xml',
        'SignUp_Signin.xml',
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
            xml = xml.replace(/{AUTHORIZATION_SERVICE_URL}/g,                    config.authorizatioServiceUrl)

            await graph.updateTrustFrameworkPolicy(policyId, xml)
    };
}

module.exports = {
    configure
}