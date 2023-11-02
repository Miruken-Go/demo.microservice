const config      = require('./config');
const axios       = require('axios').default;
const querystring = require('querystring');

let _token  = undefined

const APP_ID = "00000003-0000-0000-c000-000000000000"

async function getToken() {
    if (_token) return _token;

    const uri=`https://login.microsoftonline.com/${config.b2cName}.onmicrosoft.com/oauth2/v2.0/token`
    const result = await axios.post(uri, querystring.stringify({
        client_id:     config.b2cDeploymentPipelineClientId,
        scope:         'https://graph.microsoft.com/.default',
        client_secret: config.secrets.b2cDeploymentPipelineClientSecret,
        grant_type:    'client_credentials'
    }));
    console.log('Retrieved token')
    _token = result.data.access_token;
    return _token;
}

function logError(error){
    if(error.response){
        console.log(`status: ${error.response.status}`)
        console.log(`error.response.data: ${JSON.stringify(error.response.data)}`)
    }
}

async function get(endpoint, version) {
    vs = version || 'v1.0'
    const uri = `https://graph.microsoft.com/${vs}${endpoint}`
    console.log(`Getting: ${uri}`)

    const options = {
        headers: {
            Authorization: `Bearer ${await getToken()}`
        }
    };
    var result = await axios.get(uri, options)
        .catch(function (error) {
            console.log(`Failed to Get: ${uri}`)  
            logError(error)
            throw error
         });
    return result
}

async function post(endpoint, json, version) {
    vs = version || 'v1.0'
    const uri = `https://graph.microsoft.com/${vs}${endpoint}`
    console.log(`Posting: ${uri}`)

    const options = {
        headers: {
            Authorization: `Bearer ${await getToken()}`,
            "Content-Type": "application/json"
        }
    };
    var result = await axios.post(uri, json, options)
        .catch(function (error) {
            console.log(`Failed to Post: ${uri}`)  
            logError(error)
            throw error
         });
    return result
}

async function patch(endpoint, json, version) {
    vs = version || 'v1.0'
    const uri = `https://graph.microsoft.com/${vs}${endpoint}`
    console.log(`Patching: ${uri}`)

    const options = {
        headers: {
            Authorization: `Bearer ${await getToken()}`,
            "Content-Type": "application/json"
        }
    };
    var result = await axios.patch(uri, json, options)
        .catch(function (error) {
            console.log(`Failed to Patch: ${uri}`)  
            logError(error)
            throw error
         });

    return result;
}

module.exports = {
    get,
    post,
    patch,
    APP_ID
}
