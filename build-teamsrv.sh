docker run -v $(pwd):/go/src --workdir=/go/src/teamsrv golang:1.20 go test ./...
if [[ $? -gt 0 ]]; then 
  echo "Failed to build and test"; 
  exit 1; 
fi

TAG="teamsrv/v$(date +%s)"
if [[ $(git tag -l "$TAG") ]];
    then
        echo "Tag already created"
    else
        echo "Tagging the release"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a $TAG -m "Tagged by build pipeline"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin $TAG
fi;