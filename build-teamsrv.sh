NAME=teamsrv
VERSION="v$(date +%s)"
IMAGE_NAME="${NAME}shared.azurecr.io/${NAME}:$VERSION"; echo $IMAGE_NAME

echo "VERSION: $VERSION"
echo "IMAGE_NAME: $IMAGE_NAME"

docker build --build-arg app_version=$VERSION -t $IMAGE_NAME $NAME
if [[ $? -gt 0 ]]; then 
  echo "Failed to build docker image"; 
  exit 1; 
fi

az login --service-principal --username $DEPLOYMENT_PIPELINE_CLIENT_ID  --password $DEPLOYMENT_PIPELINE_CLIENT_SECRET --tenant $TENANT_ID
az acr login -n ${NAME}shared
docker push $IMAGE_NAME
if [[ $? -gt 0 ]]; then 
  echo "Failed to push docker image"; 
  exit 1; 
fi

TAG="$NAME/$VERSION"
echo "TAG: $TAG"
if [[ $(git tag -l "$TAG") ]];
    then
        echo "Tag already created"
    else
        echo "Tagging the release"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" tag -a $TAG -m "Tagged by build pipeline"
        git -c "user.name=buildpipeline" -c "user.email=mirukenjs@gmail.com" push origin $TAG
fi;

gh workflow run deploy-teamsrv.yml                                  \
  -f env=dev                                                        \
  -f instance=ci                                                    \
  -f tag=$VERSION                                                   \