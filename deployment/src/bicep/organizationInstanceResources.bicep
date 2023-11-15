@secure()
param containerRepositoryPassword    string 

param prefix                         string
param location                       string 
param containerRepositoryName        string 
param keyVaultName                   string
param keyVaultResourceGroup          string
param applications {
    name:             string 
    containerAppName: string 
    secrets:          string[]
}[] = [
  {
    name:             'adb2c-auth-srv'
    containerAppName: 'adb2c-auth-srv-dev'
    secrets:          []
  }
]

/////////////////////////////////////////////////////////////////////////////////////
// Container Apps
/////////////////////////////////////////////////////////////////////////////////////

resource containerAppsEnvironment 'Microsoft.App/managedEnvironments@2022-10-01' = {
  name: '${prefix}-cae'
  location: location
  sku: {
    name: 'Consumption'
  }
  properties: {
    zoneRedundant: false
    customDomainConfiguration: {}
  }
}

module containerApps 'modules/containerApp.bicep' = [for app in applications: {
  name: app.containerAppName
  params: {
    containerAppsEnvironmentId:  containerAppsEnvironment.id
    appName:                     app.name
    containerAppName:            app.containerAppName    
    prefix:                      prefix 
    location:                    location
    containerRepositoryName:     containerRepositoryName
    containerRepositoryPassword: containerRepositoryPassword
    keyVaultName:                keyVaultName
    keyVaultResourceGroup:       keyVaultResourceGroup
    secrets:                     app.secrets
  }
}]

output containerAppUrls array = [for (app, index) in applications: { 
  app: app.containerAppName 
  url: containerApps[index].outputs.containerApp.properties.configuration.ingress.fqdn
}]
