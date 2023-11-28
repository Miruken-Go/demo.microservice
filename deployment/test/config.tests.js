const { expect }  = require('chai');
const { inspect } = require('node:util');


const {
    Organization,
    Domain,
    Application,
    ResourceGroups,
    B2C
} = require('../src/infrastructure/config')

describe('ResourceGroups', function () { 
    it('exitsts', function () { 
        expect(ResourceGroups).to.exist
    })
    describe('validation', function () {
        it('name is required', function () {
            expect(()=>{new ResourceGroups({})}).to.throw('name required')
        })
        it('env is required', function () {
            expect(()=>{new ResourceGroups({name: 'foo'})}).to.throw('env required')
        })
    })

    describe('with instance', function () {
        const resourceGroups = new ResourceGroups({
            name:     'majorleaguemiruken',
            env:      'dev',
            instance: 'ci'
        })

        it('global', function () {
            expect(resourceGroups.global).to.equal('majorleaguemiruken-global')
        })

        it('common', function () {
            expect(resourceGroups.common).to.equal('majorleaguemiruken-dev-common')
        })

        it('manual', function () {
            expect(resourceGroups.manual).to.equal('majorleaguemiruken-dev-manual')
        })

        it('stable', function () {
            expect(resourceGroups.stable).to.equal('majorleaguemiruken-dev')
        })

        it('instance', function () {
            expect(resourceGroups.instance).to.equal('majorleaguemiruken-dev-ci')
        })
    })
    describe('without instance', function () {
        const resourceGroups = new ResourceGroups({
            name:     'majorleaguemiruken',
            env:      'dev'
        })

        it('global', function () {
            expect(resourceGroups.global).to.equal('majorleaguemiruken-global')
        })

        it('common', function () {
            expect(resourceGroups.common).to.equal('majorleaguemiruken-dev-common')
        })

        it('manual', function () {
            expect(resourceGroups.manual).to.equal('majorleaguemiruken-dev-manual')
        })

        it('stable', function () {
            expect(resourceGroups.stable).to.equal('majorleaguemiruken-dev')
        })

        it('instance', function () {
            expect(resourceGroups.instance).to.equal('majorleaguemiruken-dev')
        })
    })
})

describe('B2C', function () { 
    it('exitsts', function () { 
        expect(B2C).to.exist
    })
    describe('validation', function () {
        it('name is required', function () {
            expect(()=>{new B2C({})}).to.throw('name required')
        })
        it('env is required', function () {
            expect(()=>{new B2C({name: 'foo'})}).to.throw('env required')
        })
    })

    describe('b2c', function(){

        let b2c = new B2C({
            name: 'Major-League-Miruken',
            env:  'dev'
        })

        it('makes names lowercase and removes special characters', function () {
            expect(b2c.name).to.equal('majorleaguemirukenidentitydev')
        })
        it('b2cDisplayName', function () {
            expect(b2c.displayName).to.equal('majorleaguemiruken identity dev')
        })
        it('b2cDomainName', function () {
            expect(b2c.domainName).to.equal('majorleaguemirukenidentitydev.onmicrosoft.com')
        })
        it('openIdConfigurationUrl', function () {
            expect(b2c.openIdConfigurationUrl).to.equal('https://majorleaguemirukenidentitydev.b2clogin.com/majorleaguemirukenidentitydev.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN')
        })
    })
})

describe('Organization', function () {
    let org = new Organization({
        name: 'majorLeagueMiruken',
        env:  'dev'
    })
    it('exitsts', function () { 
        expect(Organization).to.exist
    })

    it('has a name', function () {
        expect(org.name).to.equal('majorleaguemiruken')
    })

    it('name is required', function () {
        expect(() => { new Organization({})}).to.throw('name required')
    })

    it('resourceGroups', function () {
        expect(org.resourceGroups).to.exist
    })

    it('globalResourceGroup', function () {
        expect(org.resourceGroups.global).to.equal('majorleaguemiruken-global')
    })

    it('b2c', function () {
        expect(org.b2c).to.exist
    })

    it('keyVaultName', function () {
        expect(org.keyVaultName).to.be.equal('majorleaguemiruken-dev-common-keyvault')
    })

    describe('containerRepository', function(){
        let org = new Organization({
            name: 'Major-League-Miruken',
            env:  'dev'
        })

        it('containerRepositoryName', function () {
            expect(org.containerRepositoryName).to.equal('majorleaguemirukenglobal')
        })
        it('containerRepositoryName validates the max length', function () {
            expect(()=>{
                new Organization({
                    name: '123456789012345678901234567',
                    env:  'dev',
                })
            }).to.throw("containerRepositoryName cannot be longer than 32 characters")
        })
    })
})

describe('Domain', function () {

    const env = 'dev'

    let org = new Organization({
        name: 'Major-League-Miruken',
        env:  env
    })

    let domain = new Domain({
        name:         'billing',
        organization: org,
        env:          env,
        applications: [
            {name: 'app1'}
        ]
    })

    it('exitsts', function () { 
        expect(Domain).to.exist
    })

    it('has a name', function () {
        expect(domain.name).to.equal('billing')
    })

    it('name is required', function () {
        expect(() => { new Domain({})}).to.throw('name required')
    })

    it('has a reference to the organization', function () {
        expect(domain.organization).to.equal(org)
    })

    it('has array of applications', function () {
        expect(domain.applications.length).to.equal(1)
    })
})

describe('Application', function () {
    it('exitsts', function () { 
        expect(Application).to.exist
    })
})

describe('Instantiating Organization', function () {
    const org = new Organization({
        name:     'MajorLeagueMiruken',
        location: 'CentralUs',
        env:      'dev',
        instance: 'ci',
        domains: [
            {
                name: 'billing', 
                applications: [
                    {
                        name: 'billingui',  
                        ui:   true
                    },
                    {
                        name: 'billingsrv', 
                        ui:   true, 
                        api:  true
                    },
                ]
            },
            {
                name: 'league', 
                applications: [
                    {
                        name: 'majorleaguemiruken', 
                        ui:   true
                    },
                    {
                        name: 'tournaments',
                        ui:   true
                    },
                    {
                        name: 'teamsrv',            
                        ui:   true, 
                        api:  true
                    },
                    {
                        name: 'schedulesrv',        
                        ui:   true, 
                        api:  true
                    },
                ]
            },
        ],
    })
    it('creates domain', function () { 
        console.log(inspect(org, { depth: null }))
        expect(org.domains.length).to.be.equal(2)
        expect(org.domains[0].instance).to.be.equal('ci')
        expect(org.domains[0].applications.length).to.be.equal(2)
        expect(org.domains[0].applications[0].instance).to.be.equal('ci')
        expect(org.domains[1].instance).to.be.equal('ci')
        expect(org.domains[1].applications.length).to.be.equal(4)
    })
})

