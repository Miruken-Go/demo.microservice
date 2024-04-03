import { 
    Domain,
    Opts,
    ContainerRepositoryResource,
    B2CResource,
    KeyVaultResource,
 } from 'ci.cd'

const env              = process.env.env
const instance         = process.env.instance
const location         = 'CentralUs'
const gitRepositoryUrl = 'https://github.com/Miruken-Go/demo.microservice'

const orgOpts: Opts = {
    name: 'MajorLeagueMiruken',
    env,
    instance,
}

export const organization = new Domain({
    ...orgOpts,
    location,
    gitRepositoryUrl,
    bootstrapUsers: [
        'provenstyle.testing@gmail.com',
        'cneuwirt@gmail.com',
    ],
    resources: {
        b2c:                 new B2CResource(orgOpts),
        containerRepository: new ContainerRepositoryResource(orgOpts),
        keyVault:            new KeyVaultResource(orgOpts),
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
            env,
            instance,
            location,
            gitRepositoryUrl,
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
            env,
            instance,
            location,
            gitRepositoryUrl,
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
