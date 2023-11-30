
export class ResourceGroups {

    constructor (opts) {
        if (!opts.name) throw new Error("name required")

        this.name    = opts.name.toLowerCase()
        this.env     = opts.env
        this.envInst = opts.instance
    }

    requireEnv () {
        if (!this.env) throw new Error("env required")
    }

    get global () {
        return `${this.name}-global`
    }

    get common () {
        this.requireEnv()
        return `${this.name}-${this.env}-common`
    }

    get manual () {
        this.requireEnv()
        return `${this.name}-${this.env}-manual`
    }

    get stable () {
        this.requireEnv()
        return `${this.name}-${this.env}`
    }

    get instance () {
        this.requireEnv()
        return (this.envInst) 
            ? `${this.name}-${this.env}-${this.envInst}`
            : `${this.name}-${this.env}`
    }
}

export class B2C {

    constructor (opts) {
        if (!opts.name) throw new Error("name required")

        this.cleanedName = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
        this.env         = opts.env
        this.profile     = opts.profile || 'B2C_1A_SIGNUP_SIGNIN'
    }

    requireEnv () {
        if (!this.env) throw new Error("env required")
    }

    get name () {
        this.requireEnv()
        return `${this.cleanedName}auth${this.env}`.toLowerCase()
    }

    get displayName () {
        this.requireEnv()
        return `${this.cleanedName} auth ${this.env}`.toLowerCase()
    }

    get domainName () {
        return `${this.name}.onmicrosoft.com`
    }

    get openIdConfigurationUrl () {
        return `https://${this.name}.b2clogin.com/${this.name}.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=${this.profile}`
    }
}

export class Application {
    name
    env
    instance
    location
    parent
    organization
    resourceGroups
    implicitFlow
    spa
    enrichApi
    scopes  = []
    secrets = []
    imageName

    constructor (opts) {
        if (!opts.name)           throw new Error("name required")
        if (!opts.location)       throw new Error("location required")
        if (!opts.organization)   throw new Error("organization required")
        if (!opts.resourceGroups) throw new Error("resourceGroups required")

        const name         = opts.name
        const organization = opts.organization

        this.name           = name 
        this.organization   = organization
        this.location       = opts.location
        this.parent         = opts.parent
        this.env            = opts.env
        this.instance       = opts.instance
        this.resourceGroups = opts.resourceGroups
        this.implicitFlow   = opts.implicitFlow || false
        this.spa            = opts.spa          || false
        this.enrichApi      = opts.enrichApi    || false
        this.scopes         = opts.scopes       || ['Groups', 'Roles', 'Entitlements']
        this.secrets        = opts.secrets      || []
        this.imageName      = `${organization.containerRepositoryName}.azurecr.io/${name}` 
    }

    get containerAppName () {
        if (!this.env) throw new Error("env required")

        const containerAppName =  (this.instance)
            ? `${this.name}-${this.env}-${this.instance}`
            : `${this.name}-${this.env}`

        if (containerAppName.length > 32)
            throw `Configuration Error - containerAppName cannot be longer than 32 characters : ${containerAppName} [${containerAppName.length}]`

        return containerAppName
    }
}

export class Domain {
    name
    env
    instance
    organization
    applications = []

    constructor (opts) {
        if (!opts.name)         throw new Error("name required")
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
                        parent:         this,
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
}

export class Organization {
    name
    env
    instance
    location
    gitRepositoryUrl
    containerRepositoryName
    resourceGroups
    b2c
    domains      = []
    applications = []

    constructor (opts) {
        if (!opts.name)     throw new Error("name required")
        if (!opts.location) throw new Error("location required")

        const name = opts.name.replace(/[^A-Za-z0-9]/g, "").toLowerCase()
        if (name.length > 19)
            throw `Configuration Error - Organization name cannot be longer than 19 characters : ${name} [${name.length}]`

        const env = opts.env
        if (env && env.length > 4)
            throw `Configuration Error - Env cannot be longer than 4 characters : ${env} [${env.length}]`

        const instance = opts.instance
        const location = opts.location

        this.name             = name
        this.env              = env 
        this.instance         = instance
        this.location         = location
        this.gitRepositoryUrl = opts.gitRepositoryUrl

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


        if(opts.applications) {
            for (const application of opts.applications) {
                this.applications.push((application instanceof Application)
                    ? application
                    : new Application({
                        parent:         this,
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

    requireEnv () {
        if (!this.env) throw new Error("env required")
    }

    get keyVaultName () {
        this.requireEnv()

        const name = `${this.name}-${this.env}` 
        if (name.length > 24)
            throw new Error(`Configuration Error - keyVaultName cannot be longer than 24 characters : ${name} [${name.length}]`)

        return name
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
