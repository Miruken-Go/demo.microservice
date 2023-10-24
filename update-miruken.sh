docker run -v $(pwd):/go/src --workdir=/go/src/teamapi golang:1.21 go get -u
if [[ $? -gt 0 ]]; then 
  echo "Failed to build and test"; 
  exit 1; 
fi

VERSION="$(docker run --rm -v "$(pwd):/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=teamapi/v)"
TAG="teamapi/v${VERSION}"
if [[ $(git tag -l "$TAG") ]];
    then
        echo "Tag already created"
    else
        echo "Tagging the release"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a $TAG -m "Tagged by build pipeline"
        #git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin $TAG
fi;




# go get -u all 

# teamapi
# go get github.com/miruken-go/miruken@v0.30.4

# teamapi
# go get github.com/miruken-go/miruken@v0.30.4
# go get github.com/miruken-go/demo.microservice/teamapi@v0.10.1

# teamsrv
# go get github.com/miruken-go/miruken@v0.30.4
# go get github.com/miruken-go/demo.microservice/teamapi@v0.10.1
# go get github.com/miruken-go/demo.microservice/team@v0.2.1
