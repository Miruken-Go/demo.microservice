import { 
    Domain,
    Opts,
    ContainerRepositoryResource,
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
    resources: {
        containerRepository: new ContainerRepositoryResource(orgOpts),
    },
    domains: [
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
