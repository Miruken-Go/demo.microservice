import * as logging     from './logging.js'
import { secrets }      from './keyvault.js'
import * as querystring from 'node:querystring'
import axios            from 'axios'

export class Graph {
    organization
    b2cDeploymentPipelineClientId
    _token  = undefined

    static APP_ID = "00000003-0000-0000-c000-000000000000"

    constructor (organization, b2cDeploymentPipelineClientId) {
        if (!organization)                  throw new Error('organization is required')
        if (!b2cDeploymentPipelineClientId) throw new Error('b2cDeploymentPipelineClientId is required')

        this.organization                  = organization
        this.b2cDeploymentPipelineClientId = b2cDeploymentPipelineClientId
    }

    async getToken() {
        if (this._token) return this._token;

        await secrets.requireSecrets([
            'b2cDeploymentPipelineClientSecret',
        ], this.organization.keyVaultName)

        logging.printEnvironmentSecrets(secrets)

        const uri=`https://login.microsoftonline.com/${this.organization.b2c.domainName}/oauth2/v2.0/token`
        const result = await axios.post(uri, querystring.stringify({
            client_id:     this.b2cDeploymentPipelineClientId,
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

        try {
            return await axios.get(uri, options)
        } catch (error) {
            console.log(`Failed to Get: ${uri}`)  
            this.logError(error)
            throw error
        }
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

        try {
            return await axios.post(uri, json, options)
        } catch (error) {
            console.log(`Failed to Post: ${uri}`)  
            this.logError(error)
            throw error
        }
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

        try {
            return await axios.patch(uri, json, options)
        } catch (error) {
            console.log(`Failed to Patch: ${uri}`)  
            this.logError(error)
            throw error
        }
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

        try {
            var result = await axios.put(uri, xml, options)
            console.log(result.status)
            return result;
        } catch (error) {
            console.log(`Failed to Update: ${uri}`)  
            this.logError(error)
            throw error
        }
    }
}
