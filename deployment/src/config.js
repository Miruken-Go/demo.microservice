import { fileURLToPath } from 'node:url'
import { dirname }       from 'node:path'
import { 
    Domain,
    B2C,
    ContainerRepository,
    KeyVault,
} from '#infrastructure/config.js'

export const configDirectory = dirname(fileURLToPath(import.meta.url))

export const organization = new Domain({
    env:              process.env.env,
    instance:         process.env.instance,
    name:             'MajorLeagueMiruken',
    location:         'CentralUs',
    gitRepositoryUrl: 'https://github.com/Miruken-Go/demo.microservice',
    bootstrapUsers: [
        'provenstyle.testing@gmail.com',
        'cneuwirt@gmail.com',
    ],
    resources: {
        b2c:                 B2C,
        containerRepository: ContainerRepository,
        keyVault:            KeyVault,
    },
    applications: [
        {
            name:      'adb2c-api-connector-srv', 
            enrichApi: true,  
            secrets: [
                'authorization-service-password',
                'cosmos-connection-string',
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
