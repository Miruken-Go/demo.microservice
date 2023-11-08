const graph    = require('./graph');

async function getApplications() {
    const result = await graph.get("/applications")
    return result.data.value
}

async function getApplicationById(id) {
    const result = await graph.get(`/applications/${id}`)
    return result.data
}

async function getApplicationByName(displayName) {
    const applications = await getApplications()
    return applications.find(a => a.displayName === displayName)
}

async function updateApplication(id, manifest) {
    console.log(`Updating existing application appId [${id}]`)
    await graph.patch(`/applications/${id}`, manifest)
    return await getApplicationById(id)
}

async function createOrUpdateApplication(manifest) {
    const displayName = manifest.displayName
    const existing = await getApplicationByName(displayName)

    let application = undefined
    if (existing) {
        application = await updateApplication(existing.id, manifest)
    } else {
        console.log(`Creating application: ${displayName}`)
        application = (await graph.post("/applications", manifest)).data
        console.log(`Created Application: ${displayName}`)
        console.log(application)
    }

    return application
}

async function addRedirectUris(id, uris) {
    const app = await getApplicationById(id)
    const redirectUris = [...app.spa.redirectUris]
    for (const uri of uris) {
        if (!redirectUris.includes(uri)) {
            redirectUris.push(uri)
        }
    }
    await updateApplication(id, {
        spa: {
            redirectUris: redirectUris
        }
    })
    return await getApplicationById(id)
}

module.exports = {
    getApplications,
    getApplicationById,
    getApplicationByName,
    updateApplication,
    createOrUpdateApplication,
    addRedirectUris,
}
