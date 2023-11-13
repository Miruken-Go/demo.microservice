
class ResourceGroups {
    global
    common
    manual
    stable 
    instance

    constructor (opts) {
        if (!opts.name) throw new Error("name required")
        if (!opts.env)  throw new Error("env required")

        const name     = opts.name.toLowerCase()
        const env      = opts.env
        const instance = opts.instance

        this.global   = `${name}-global`
        this.common   = `${name}-${env}-common`
        this.manual   = `${name}-${env}-manual`
        this.stable   = `${name}-${env}`
        this.instance = (instance) 
            ? `${name}-${env}-${instance}`
            : `${name}-${env}`
    }
}

class B2C {
    name
    displayName
    domainName
    openIdConfigurationUrl

    constructor (opts) {
        if (!opts.name) throw new Error("name required")
        if (!opts.env)  throw new Error("env required")

        const profile  = opts.profile || 'B2C_1A_SIGNUP_SIGNIN'
        const name     = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
        const env      = opts.env
        const b2cName  = `${name}identity${env}`.toLowerCase()

        this.name                   = b2cName
        this.displayName            = `${name} identity ${env}`.toLowerCase()
        this.domainName             = `${b2cName}.onmicrosoft.com`
        this.openIdConfigurationUrl = `https://${b2cName}.b2clogin.com/${b2cName}.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=${profile}`
    }
}

class Application {
    name
    organization
    domain
    api
    ui

    constructor (opts) {
        this.name         = opts.name.toLowerCase()
        this.organization = opts.organization
        this.domain       = opts.domain
        this.api          = opts.api || false
        this.ui           = opts.ui  || false
    }

    get containerAppEnvironmentName () {
        return `${domain.instancePrefix}-cae`
    }

    get containerAppName () {
        return `${domain.instancePrefix}-${this.name}`
    }

    get imageName () { 
        return `${this.organization.containerRepositoryName}.azurecr.io/${appName}` 
    }
}

class Domain {
    name
    env
    instance
    organization
    apps = []

    constructor (opts) {
        if (!opts.name) throw new Error("name required")
        if (!opts.env)  throw new Error("env required")

        this.name         = opts.name
        this.env          = opts.env
        this.organization = opts.organization
        this.apps         = opts.apps

        this.resourceGroups = new ResourceGroups({
            name:     opts.name,
            env:      opts.env,
            instance: opts.instance,
        })
    }

    get commonPrefix () {
        return `${this.name}-${this.env}-common`
    }
    get instancePrefix () {
        return (instance) 
        ? `${this.name}-${this.env}-${instance}`
        : `${this.name}-${this.env}`
    }

    get keyVaultName () {
        return `${this.commonPrefix}-keyvault` 
    }
}


class Organization {
    name
    env
    instance
    containerRepositoryName
    domains = []
    resourceGroups
    b2c

    constructor (opts) {
        if (!opts.name) throw new Error("name required")
        if (!opts.env)  throw new Error("env required")

        const name = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()

        this.name     = name
        this.env      = opts.env
        this.instance = opts.instance 
        this.domains  = opts.domains

        this.containerRepositoryName = `${name}global`
        if (this.containerRepositoryName.length > 32)
            throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${this.containerRepositoryName} [${this.containerRepositoryName.length}]`

        this.resourceGroups = new ResourceGroups({
            name: name, 
            env:  opts.env
        })

        this.b2c = new B2C({
            name: name, 
            env:  opts.env
        })
    }
}

module.exports = {
    ResourceGroups,
    B2C,
    Organization,
    Domain,
    Application,
}