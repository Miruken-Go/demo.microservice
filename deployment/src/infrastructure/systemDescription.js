class Application {
    name
    api
    ui

    constructor (name, api, ui) {
        this.name  = name
        this.api   = api
        this.ui    = ui
    }
}

class Domain {
    name
    apps = []

    constructor (name, apps) {
        this.name  = name
        this.apps  = apps
    }
}

class Organization {
    name
    domains = []

    constructor (name, domains) {
        this.name     = name
        this.domains  = domains
    }
}


class ApplicationType {
    static apiWithOpenApiUI = 'apiWithOpenApiUI'
}

module.exports = {
    ApplicationType
}