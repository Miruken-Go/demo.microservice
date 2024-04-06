@secure()
param containerRepositoryPassword    string 

param prefix                         string
param location                       string 
param containerRepositoryName        string 
param applications {
    name:             string 
    containerAppName: string 
    imageTag:         string
}[]
param tags {
    organization: string
    domain:       string
    env:          string
    instance:     string
}

/////////////////////////////////////////////////////////////////////////////////////
// Container Apps
/////////////////////////////////////////////////////////////////////////////////////

resource containerAppsEnvironment 'Microsoft.App/managedEnvironments@2022-10-01' = {
  name: '${prefix}-cae'
  location: location
  tags: tags
  sku: {
    name: 'Consumption'
  }
  properties: {
    zoneRedundant: false
    customDomainConfiguration: {}
  }
}

resource containerApps 'Microsoft.App/containerApps@2023-05-01' = [for app in applications: {
  name:     app.containerAppName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties:{
    managedEnvironmentId: containerAppsEnvironment.id
    configuration: {
      activeRevisionsMode: 'Multiple'
      ingress: {
        targetPort: 8080
        external:   true
      }
      secrets: [
        {
          name: 'acr-password'
          value: containerRepositoryPassword
        } 
      ]
      registries: [
        {
          passwordSecretRef: 'acr-password'
          username:          containerRepositoryName
          server:            '${containerRepositoryName}.azurecr.io'
        }
      ]
    }
    template: {
      containers: [
        {
          image: '${containerRepositoryName}.azurecr.io/${app.name}:${app.imageTag}' 
          name:  app.name
          env: [
            {
              name: 'RESOURCE_GROUP'
              value: resourceGroup().name
            }
            {
              name: 'APPLICATION_NAME'
              value: app.name
            }
          ]
        }
      ]
    }
  }
}]

output containerAppUrls array = [for (app, index) in applications: { 
  app: containerApps[index].name
  url: containerApps[index].properties.configuration.ingress.fqdn
}]
