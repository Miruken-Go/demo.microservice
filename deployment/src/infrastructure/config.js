
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
        const b2cName  = `${name}auth${env}`.toLowerCase()

        this.name                   = b2cName
        this.displayName            = `${name} auth ${env}`.toLowerCase()
        this.domainName             = `${b2cName}.onmicrosoft.com`
        this.openIdConfigurationUrl = `https://${b2cName}.b2clogin.com/${b2cName}.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=${profile}`
    }
}

class Application {
    name
    env
    instance
    location
    organization
    resourceGroups
    implicitFlow
    spa
    enrichApi
    scopes  = []
    secrets = []
    imageName
    containerAppName

    constructor (opts) {
        if (!opts.name)           throw new Error("name required")
        if (!opts.env)            throw new Error("env required")
        if (!opts.location)       throw new Error("location required")
        if (!opts.organization)   throw new Error("organization required")
        if (!opts.resourceGroups) throw new Error("resourceGroups required")

        const name         = opts.name
        const env          = opts.env
        const instance     = opts.instance
        const location     = opts.location
        const organization = opts.organization

        this.name           = name 
        this.env            = env
        this.instance       = instance
        this.location       = location
        this.organization   = organization
        this.resourceGroups = opts.resourceGroups
        this.implicitFlow   = opts.implicitFlow || false
        this.spa            = opts.spa          || false
        this.enrichApi      = opts.enrichApi    || false
        this.scopes         = opts.scopes       || ['Groups', 'Roles', 'Entitlements']
        this.secrets        = opts.secrets
        this.imageName      = `${organization.containerRepositoryName}.azurecr.io/${name}` 
        this.containerAppName = (instance)
            ? `${name}-${env}-${instance}`
            : `${name}-${env}`

        if (this.containerAppName.length > 32)
            throw `Configuration Error - containerAppName cannot be longer than 32 characters : ${this.containerAppName} [${this.containerAppName.length}]`
    }


    // get containerAppEnvironmentName () {
    //     return `${domain.instancePrefix}-cae`
    // }

    // get containerAppName () {
    //     return `${domain.instancePrefix}-${this.name}`
    // }
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
        if (!opts.location)     throw new Error("location required")
        if (!opts.organization) throw new Error("organization required")

        const name     = opts.name
        const env      = opts.env
        const instance = opts.instance
        const location = opts.location

        this.name         = name
        this.env          = env
        this.instance     = instance
        this.location     = location
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
                        domain:         this,
                        organization:   this.organization,
                        env:            env, 
                        instance:       instance,
                        location:       location,
                        resourceGroups: this.resourceGroups,
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

}


class Organization {
    name
    env
    instance
    location
    containerRepositoryName
    resourceGroups
    b2c
    keyVaultName
    domains      = []
    applications = []

    constructor (opts) {
        if (!opts.name)     throw new Error("name required")
        if (!opts.env)      throw new Error("env required")
        if (!opts.location) throw new Error("location required")

        const name = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
        if (name.length > 19)
            throw `Configuration Error - Organization name cannot be longer than 19 characters : ${name} [${name.length}]`

        const env = opts.env
        if (env.length > 4)
            throw `Configuration Error - Env cannot be longer than 4 characters : ${env} [${env.length}]`

        const instance = opts.instance
        const location = opts.location

        this.name     = name
        this.env      = env 
        this.instance = instance
        this.location = location

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

        this.keyVaultName = `${name}-${env}` 
        if (this.keyVaultName.length > 24)
            throw `Configuration Error - keyVaultName cannot be longer than 24 characters : ${this.keyVaultName} [${this.keyVaultName.length}]`

        if(opts.applications) {
            for (const application of opts.applications) {
                this.applications.push((application instanceof Application)
                    ? application
                    : new Application({
                        organization:   this,
                        env:            env, 
                        instance:       instance,
                        location:       location,
                        resourceGroups: this.resourceGroups,
                        ...application,
                    }))
            }
        }

        if(opts.domains) {
            for (const domain of opts.domains) {
                this.domains.push((domain instanceof Domain)
                    ? domain
                    : new Domain({
                        organization: this,
                        env:          env,
                        instance:     instance,
                        location:     location, 
                        ...domain,
                    }))
            }
        }
    }

    get enrichApiApplication () {
        let application  = this.applications.find(a => a.enrichApi)
        if (application) return application

        throw new Error(`No application defined in organization where enrichApi = true`)
    }

    getApplicationByName(name) {
        let application = this.applications.find(a => a.name === name)
        if (application) return application

        for (const domain of this.domains) {
            application = domain.applications.find(a => a.name === name)
            if (application) return application
        }

        throw new Error(`Application with name ${name} not found`)
    }
}

module.exports = {
    ResourceGroups,
    B2C,
    Organization,
    Domain,
    Application,
}