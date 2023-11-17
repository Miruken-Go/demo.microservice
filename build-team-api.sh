docker run -v $(pwd):/go/src --workdir=/go/src/team-api golang:1.21.3-alpine3.18 go test ./...
if [[ $? -gt 0 ]]; then 
  echo "Failed to build and test"; 
  exit 1; 
fi

VERSION="$(docker run --rm -v "$(pwd):/repo" gittools/gitversion:5.12.0-alpine.3.14-6.0 /repo /showvariable SemVer /overrideconfig tag-prefix=team-api/v)"
TAG="team-api/v${VERSION}"
if [[ $(git tag -l "$TAG") ]];
    then
        echo "Tag already created"
    else
        echo "Tagging the release"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a $TAG -m "Tagged by build pipeline"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin $TAG
fi;