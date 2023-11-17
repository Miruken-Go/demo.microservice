const logging       = require('./logging');
const { secrets }   = require('./keyvault');
const { variables } = require('./envVariables')
const querystring   = require('querystring');
const axios         = require('axios').default;

variables.requireEnvFileVariables([
    'b2cDeploymentPipelineClientId'
])

class Graph {
    organization
    _token  = undefined

    static APP_ID = "00000003-0000-0000-c000-000000000000"

    constructor (organization) {
        this.organization = organization
    }

    async getToken() {
        if (this._token) return this._token;

        await secrets.requireSecrets([
            'b2cDeploymentPipelineClientSecret',
        ], this.organization.keyVaultName)

        logging.printEnvironmentSecrets(secrets)

        const uri=`https://login.microsoftonline.com/${this.organization.b2c.domainName}/oauth2/v2.0/token`
        const result = await axios.post(uri, querystring.stringify({
            client_id:     variables.b2cDeploymentPipelineClientId,
            scope:         'https://graph.microsoft.com/.default',
            client_secret: secrets.b2cDeploymentPipelineClientSecret,
            grant_type:    'client_credentials'
        }));
        console.log('Retrieved token')
        this._token = result.data.access_token;
        return this._token;
    }

    logError(error){
        if(error.response){
            console.log(`status: ${error.response.status}`)
            console.log(`error.response.data: ${JSON.stringify(error.response.data)}`)
        }
    }

    async get(endpoint, version) {
        const vs = version || 'v1.0'
        const uri = `https://graph.microsoft.com/${vs}${endpoint}`
        console.log(`Getting: ${uri}`)

        const options = {
            headers: {
                Authorization: `Bearer ${await this.getToken()}`
            }
        };
        var result = await axios.get(uri, options)
            .catch(function (error) {
                console.log(`Failed to Get: ${uri}`)  
                this.logError(error)
                throw error
            });
        return result
    }

    async post(endpoint, json, version) {
        const vs = version || 'v1.0'
        const uri = `https://graph.microsoft.com/${vs}${endpoint}`
        console.log(`Posting: ${uri}`)

        const options = {
            headers: {
                Authorization: `Bearer ${await this.getToken()}`,
                "Content-Type": "application/json"
            }
        };
        var result = await axios.post(uri, json, options)
            .catch(function (error) {
                console.log(`Failed to Post: ${uri}`)  
                this.logError(error)
                throw error
            });
        return result
    }

    async patch(endpoint, json, version) {
        const vs = version || 'v1.0'
        const uri = `https://graph.microsoft.com/${vs}${endpoint}`
        console.log(`Patching: ${uri}`)

        const options = {
            headers: {
                Authorization: `Bearer ${await this.getToken()}`,
                "Content-Type": "application/json"
            }
        };
        var result = await axios.patch(uri, json, options)
            .catch(function (error) {
                console.log(`Failed to Patch: ${uri}`)  
                this.logError(error)
                throw error
            });

        return result;
    }

    //https://learn.microsoft.com/en-us/graph/api/trustframework-put-trustframeworkpolicy?view=graph-rest-beta
    async updateTrustFrameworkPolicy(policyId, xml) {
        const uri = `https://graph.microsoft.com/beta/trustFramework/policies/${policyId}/$value`
        console.log(`Putting: ${uri}`)

        const options = {
            headers: {
                Authorization: `Bearer ${await this.getToken()}`,
                "Content-Type": "application/xml"
            }
        };
        var result = await axios.put(uri, xml, options)
            .catch(function (error) {
                console.log(`Failed to Update: ${uri}`)  
                this.logError(error)
                throw error
            });

        console.log(result.status)

        return result;
    }
}

module.exports = {
    Graph
}
