param prefix                      string
param appName                     string
param location                    string 
param containerRepositoryName     string 
@secure()
param containerRepositoryPassword string 


/////////////////////////////////////////////////////////////////////////////////////
// Azure Monitor
/////////////////////////////////////////////////////////////////////////////////////

// resource logAnalyticsWorkspace'microsoft.operationalinsights/workspaces@2021-06-01' = {
//   name: '${prefix}-log-analytics-workspace'
//   location: location
//   properties: {
//     sku: {
//       name: 'PerGB2018'
//     }
//     retentionInDays: 30
//     features: {
//       enableLogAccessUsingOnlyResourcePermissions: true
//     }
//   }
// }

/////////////////////////////////////////////////////////////////////////////////////
// Container Apps
/////////////////////////////////////////////////////////////////////////////////////
resource containerAppsEnvironment 'Microsoft.App/managedEnvironments@2022-10-01' = {
  name: '${prefix}-CAE'
  location: location
  sku: {
    name: 'Consumption'
  }
  properties: {
    zoneRedundant: false
    customDomainConfiguration: {}
  }
}

resource containerApp 'Microsoft.App/containerApps@2022-03-01' ={
  name: prefix
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
        external: true
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
          username: containerRepositoryName
          server: '${containerRepositoryName}.azurecr.io'
        }
      ]
    }
    template: {
      containers: [
        {
          image: '${containerRepositoryName}.azurecr.io/${appName}:initial' 
          name:  appName
        }
      ]
    }
  }
}
