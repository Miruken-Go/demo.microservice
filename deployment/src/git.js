const bash    = require('./bash')
const config  = require('./config');
const logging = require('./logging');

console.log("Configuring git")
config.requiredSecrets(['ghToken'])
bash.execute(`
    git config --global --add safe.directory $(pwd)
    git config --global url."https://api:$ghToken@github.com/".insteadOf "https://github.com/"
    git config --global url."https://ssh:$ghToken@github.com/".insteadOf "ssh://git@github.com/"
    git config --global url."https://git:$ghToken@github.com/".insteadOf "git@github.com:"
`)

async function tagAndPush(tag) { 
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
            git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a ${tag} -m "Tagged by build pipeline"
            git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin ${tag}
        `)
    }
}

async function anyChanges() { 
    const status = await bash.execute(`
        git status
    `)
    return status.substring('Changes not staged for commit');
}

async function commitAll(message) { 
    logging.header("Commiting Changes")

    await bash.execute(`
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" commit -am "${message}"
    `)
}

async function push() { 
    logging.header("Pushing branch")
    await bash.execute(`
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin
    `)
}

module.exports = {
    anyChanges,
    commitAll,
    push,
    tagAndPush
}
