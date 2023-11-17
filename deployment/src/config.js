const { Organization } = require('./infrastructure/config')

const env = process.env.env
if (!env) throw "Environment variable required: [env]"

const instance = process.env.instance

const org = new Organization({
    name:     'MajorLeagueMiruken',
    location: 'CentralUs',
    env:      env,
    instance: instance,
    applications: [
        {
            name:        'adb2c-auth-srv', 
            ui:          true, 
            api:         true,
            isEnrichApi: true,  
            secrets: [
                'authorization-service-password'
            ]
        },
    ],
    domains: [
        {
            name: 'billing', 
            applications: [
                {
                    name: 'billing-ui',  
                    ui:   true
                },
                {
                    name: 'billing-srv', 
                    ui:   true, 
                    api:  true
                },
            ]
        },
        {
            name: 'league', 
            applications: [
                {
                    name: 'major-league-miruken-ui', 
                    ui:   true
                },
                {
                    name: 'tournaments-ui',
                    ui:   true
                },
                {
                    name: 'teamsrv',            
                    ui:   true, 
                    api:  true
                },
                {
                    name: 'schedule-srv',        
                    ui:   true, 
                    api:  true
                },
            ]
        },
    ],
})

module.exports = {
    organization: org
}
