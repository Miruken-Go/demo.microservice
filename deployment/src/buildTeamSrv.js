const az      = require('./az');
const bash    = require('./bash')
const logging = require('./logging');
const config  = require('./config');

async function main() {
    try {
        console.log("Building teamsrv")
        config.requiredNonSecrets([
        ])
        logging.printConfiguration(config)
        await az.login()

        // if [[ $(git tag -l "$TAG") ]];
        //     then
        //         echo "Tag already created"
        //     else
        //         echo "Tagging the release"
        //         git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a $TAG -m "Tagged by build pipeline"
        //         git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin $TAG
        // fi;

        // gh workflow run deploy-teamsrv.yml \
        // -f env=dev                       \
        // -f instance=ci                   \
        // -f tag=$VERSION                  \

        //Create the new revision
        const name      = config.appName
        const version   = `v${Math.floor(Date.now()/1000)}`
        const imageName = `${config.imageName}:${version}`
        const tag       = `${config.appName}/${version}`

        console.log(`version:   ${version}`);
        console.log(`imageName: ${imageName}`);
        console.log(`tag:       ${tag}`);

        await bash.execute(`
            docker build --build-arg app_version=${version} -t ${imageName} /build/demo.microservice/teamsrv
        `)
        await bash.execute(`
            az acr login -n ${config.containerRepositoryName}
        `)
        await bash.execute(`
            docker push ${imageName}
        `)

        await bash.execute(`
            pwd; ls -la
        `)

        const existingTag = await bash.execute(`
            git tag -l ${tag}
        `)
        if (existingTag === tag) {
            console.log("Tag already created")
        } else {
            console.log("Tagging the release")
            await bash.execute(`
                git config --global url."https://api:$GH_TOKEN@github.com/".insteadOf "https://github.com/"
                git config --global url."https://ssh:$GH_TOKEN@github.com/".insteadOf "ssh://git@github.com/"
                git config --global url."https://git:$GH_TOKEN@github.com/".insteadOf "git@github.com:"
            `)
            await bash.execute(`
                git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a ${tag} -m "Tagged by build pipeline"
            `)
            await bash.execute(`
                git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin ${tag}
            `)
        }

        console.log("Built teamsrv")
    } catch (error) {
        process.exitCode = 1
        console.log(error)
        console.log("Deployment Failed")
    }
}

main()
