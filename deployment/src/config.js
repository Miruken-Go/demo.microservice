import { fileURLToPath } from 'url';
import { dirname }       from 'node:path'
import { Organization }  from '#infrastructure/config.js'

export const configDirectory = dirname(fileURLToPath(import.meta.url))

export const organization = new Organization({
    env:              process.env.env,
    instance:         process.env.instance,
    name:             'MajorLeagueMiruken',
    location:         'CentralUs',
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
