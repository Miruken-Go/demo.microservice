
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
    env
    instance
    organization
    domain
    api
    ui

    constructor (opts) {
        if (!opts.name)         throw new Error("name required")
        if (!opts.env)          throw new Error("env required")
        if (!opts.organization) throw new Error("organization required")
        if (!opts.domain)       throw new Error("domain required")

        const name     = opts.name
        const env      = opts.env
        const instance = opts.instance

        this.name         = name 
        this.env          = env
        this.instance     = instance
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
    applications = []

    constructor (opts) {
        if (!opts.name)         throw new Error("name required")
        if (!opts.env)          throw new Error("env  required")
        if (!opts.organization) throw new Error("organization required")

        const name     = opts.name
        const env      = opts.env
        const instance = opts.instance

        this.name         = name
        this.env          = env
        this.instance     = instance
        this.organization = opts.organization

        this.resourceGroups = new ResourceGroups({
            name:     name,
            env:      env,
            instance: instance,
        })

        if(opts.applications) {
            for (const application of opts.applications) {
                this.applications.push((application instanceof Application)
                    ? application
                    : new Application({
                        domain:       this,
                        organization: this.organization,
                        env:          env, 
                        instance:     instance,
                        ...application,
                    }))
            }
        }
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

        const name     = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
        const env      = opts.env
        const instance = opts.instance

        this.name     = name
        this.env      = env 
        this.instance = instance

        this.containerRepositoryName = `${name}global`
        if (this.containerRepositoryName.length > 32)
            throw `Configuration Error - containerRepositoryName cannot be longer than 32 characters : ${this.containerRepositoryName} [${this.containerRepositoryName.length}]`

        this.resourceGroups = new ResourceGroups({
            name:     name, 
            env:      env,
            instance: instance,
        })

        this.b2c = new B2C({
            name: name, 
            env:  env
        })

        if(opts.domains) {
            for (const domain of opts.domains) {
                this.domains.push((domain instanceof Domain)
                    ? domain
                    : new Domain({
                        organization: this,
                        env:          env,
                        instance:     instance,
                        ...domain,
                    }))
            }
        }
    }
}

module.exports = {
    ResourceGroups,
    B2C,
    Organization,
    Domain,
    Application,
}