import * as bash from './bash.js'

export async function getModuleVersion(folder, module) { 
    return await bash.execute(`
        cd ${folder}
        go list -m all | grep ${module} | awk '{print $2}' \
    `)
}
