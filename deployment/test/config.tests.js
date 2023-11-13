const {expect} = require('chai');

const {
    Organization,
    Domain,
    Application
} = require('../src/config')

describe('Organization', function () {
    let org = new Organization({
        name: 'majorLeagueMiruken'
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

    it('globalPrefix', function () {
        expect(org.globalPrefix).to.equal('majorleaguemiruken')
    })

    it('globalResourceGroup', function () {
        expect(org.globalResourceGroup).to.equal('majorleaguemiruken-global')
    })

    describe('b2c', function(){
        let org = new Organization({
            name: 'Major-League-Miruken'
        })

        it('makes names lowercase and removes special characters', function () {
            expect(org.b2cName).to.equal('majorleaguemirukenidentitydev')
        })
        it('b2cDisplayName', function () {
            expect(org.b2cDisplayName).to.equal('majorleaguemiruken identity dev')
        })
        it('b2cDomainName', function () {
            expect(org.b2cDomainName).to.equal('majorleaguemirukenidentitydev.onmicrosoft.com')
        })
        it('openIdConfigurationUrl', function () {
            expect(org.openIdConfigurationUrl).to.equal('https://majorleaguemirukenidentitydev.b2clogin.com/majorleaguemirukenidentitydev.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN')
        })
        it('containerRepositoryName', function () {
            expect(org.containerRepositoryName).to.equal('majorleaguemirukenglobal')
        })
        it('containerRepositoryName validates the max length', function () {
            let org = new Organization({
                name: '123456789012345678901234567'
            })
            expect(()=>{org.containerRepositoryName}).to.throw("containerRepositoryName cannot be longer than 32 characters")
        })
    })
})

describe('Domain', function () {
    it('exitsts', function () { 
        expect(Domain).to.exist
    })
})

describe('Application', function () {
    it('exitsts', function () { 
        expect(Application).to.exist
    })
})