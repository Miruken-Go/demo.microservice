const { Organization } = require('./infrastructure/config')

const env = process.env.env
//if (!env) throw "Environment variable required: [env]"

const instance = process.env.instance

const org = new Organization({
    name:             'MajorLeagueMiruken',
    location:         'CentralUs',
    env:              env,
    instance:         instance,
    gitRepositoryUrl: 'https://github.com/Miruken-Go/demo.microservice',
    applications: [
        {
            name:      'adb2c-api-connector-srv', 
            enrichApi: true,  
            secrets: [
                'authorization-service-password'
            ]
        },
        {
            name:         'adb2c-auth-srv', 
            implicitFlow: true,
            spa:          true,
        },
    ],
    domains: [
        {
            name: 'billing', 
            applications: [
                {
                    name:         'billing-srv', 
                    implicitFlow: true,
                    spa:          true,
                },
            ]
        },
        {
            name: 'league', 
            applications: [
                {
                    name: 'league-srv', 
                    implicitFlow: true,
                    spa:          true,
                },
                {
                    name: 'tournaments-srv',
                    implicitFlow: true,
                    spa:          true,
                },
                {
                    name: 'team-srv',            
                    implicitFlow: true,
                    spa:          true,
                },
                {
                    name: 'schedule-srv',        
                    implicitFlow: true,
                    spa:          true,
                },
            ]
        },
    ],
})

module.exports = {
    organization: org
}
