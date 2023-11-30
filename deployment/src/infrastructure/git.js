import * as bash    from './bash.js'
import * as logging from './logging.js'
import { secrets }  from './envSecrets.js'

secrets.require([
   'ghToken' 
])

console.log("Configuring git")
bash.execute(`
    git config --global --add safe.directory $(pwd)
    git config --global user.email "mirukenjs@gmail.com"
    git config --global user.name "buildpipeline"
    git config --global url."https://api:${secrets.ghToken}@github.com/".insteadOf "https://github.com/"
    git config --global url."https://ssh:${secrets.ghToken}@github.com/".insteadOf "ssh://git@github.com/"
    git config --global url."https://git:${secrets.ghToken}@github.com/".insteadOf "git@github.com:"
`)

export async function tagAndPush(tag) { 
    logging.header("Tagging the commit")

    const existingTag = await bash.execute(`
        git tag -l ${tag}
    `)

    console.log(`existingTag: [${existingTag}]`)
    console.log(`tag: [${tag}]`)

    if (existingTag === tag) {
        console.log("Tag already created")
    } else {
        console.log("Tagging the release")
        await bash.execute(`
            git tag -a ${tag} -m "Tagged by build pipeline"
            git push origin ${tag}
        `)
    }
}

export async function anyChanges() { 
    const status = await bash.execute(`
        git status
    `)
    const foundChanges = status.includes('Changes not staged for commit');
    if (foundChanges) {
        console.log("Changes found in git repo")
    } else {
        console.log("No changes found in git repo")
    }
    return foundChanges
}

export async function commitAll(message) { 
    logging.header("Commiting Changes")

    await bash.execute(`
        git commit -am "${message}"
    `)
}

export async function push() { 
    logging.header("Pushing branch")
    await bash.execute(`
        git push origin
    `)
}
