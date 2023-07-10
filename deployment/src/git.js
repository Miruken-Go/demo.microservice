const bash    = require('./bash')
const config  = require('./config');

async function tagAndPush(tag) { 
    config.requiredSecrets(['ghToken'])

    const existingTag = await bash.execute(`
        git tag -l ${tag}
    `)
    if (existingTag === tag) {
        console.log("Tag already created")
    } else {
        console.log("Tagging the release")
        await bash.execute(`
            git config --global url."https://api:$ghToken@github.com/".insteadOf "https://github.com/"
            git config --global url."https://ssh:$ghToken@github.com/".insteadOf "ssh://git@github.com/"
            git config --global url."https://git:$ghToken@github.com/".insteadOf "git@github.com:"
            git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a ${tag} -m "Tagged by build pipeline"
            git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin ${tag}
        `)
    }
}

module.exports = {
    tagAndPush
}
