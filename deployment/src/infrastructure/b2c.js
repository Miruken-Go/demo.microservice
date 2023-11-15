const config             = require('../config');
const logging            = require('./logging');
const graph              = require('./graph');
const az                 = require('./az');
const b2cAppRegistration = require('./b2cAppRegistration')
const fs                 = require('fs')
const path               = require('path')
const axios              = require('axios')

async function configureCustomPolicies () {

    logging.header("Deploying B2C Configuration")

    const identityExperienceFrameworkClient = await b2cAppRegistration.getApplicationByName('IdentityExperienceFramework')
    if(!identityExperienceFrameworkClient) throw new Error("IdentityExperienceFramework application not found. Check that the App Registration was created in B2C and check the name spelling and casing.")

    const proxyIdentityExperienceFrameworkClient = await b2cAppRegistration.getApplicationByName('ProxyIdentityExperienceFramework')
    if(!proxyIdentityExperienceFrameworkClient) throw new Error("ProxyIdentityExperienceFramework application not found. Check that the App Registration was created in B2C and check the name spelling and casing.")

    const appUrl = await az.getContainerAppUrl(config.prefix)
    if(!appUrl) throw new Error(`authorizationServiceUrl could not be calculated. The AppUrl for ${config.prefix} container app was not found. The default application environment instance needs to be deployed before common configuration can run.`)

    const authorizationServiceUrl = `https://${appUrl}/enrich/`

    //https://learn.microsoft.com/en-us/azure/active-directory-b2c/deploy-custom-policies-devops
    const customPoliciesDirectory = path.resolve(__dirname, '../custom-policies')
    const customPoliciesFileOrder = [
        'TrustFrameworkBase.xml',
        'TrustFrameworkLocalization.xml',
        'TrustFrameworkExtensions.xml',
        'SignUp_SignIn.xml',
        'ProfileEdit.xml',
        'PasswordReset.xml',
    ]
    for (const file of customPoliciesFileOrder) {
            const policyId = `B2C_1A_${path.basename(file, '.xml')}`
            const filePath = path.join(customPoliciesDirectory, file)
            let xml        = fs.readFileSync(filePath,{encoding: 'utf-8'}) 

            xml = xml.replace(/{B2C_DOMAIN_NAME}/g,                              config.b2cDomainName)
            xml = xml.replace(/{IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g,       identityExperienceFrameworkClient.appId)
            xml = xml.replace(/{PROXY_IDENTITY_EXPERIENCE_FRAMEWORK_CLIENTID}/g, proxyIdentityExperienceFrameworkClient.appId)
            xml = xml.replace(/{AUTHORIZATION_SERVICE_URL}/g,                    authorizationServiceUrl)

            await graph.updateTrustFrameworkPolicy(policyId, xml)
    };
}

async function getWellKnownOpenIdConfiguration() {
    const uri = config.openIdConfigurationUrl
    console.log(`Getting: ${uri}`)
    const result = await axios.get(uri)
        .catch(function (error) {
            console.log(`Failed to Get: ${uri}`)  
            logError(error)
            throw error
         });
   
    console.log(result.data)
    return result.data
}

module.exports = {
    configureCustomPolicies,
    getWellKnownOpenIdConfiguration
}