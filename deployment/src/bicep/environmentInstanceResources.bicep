param prefix                         string
param appName                        string
param location                       string 
param containerRepositoryName        string 
@secure()
param containerRepositoryPassword    string 
param keyVaultName                   string
param commonEnvironmentResourceGroup string


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

resource containerApp 'Microsoft.App/containerApps@2023-05-01' ={
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
        {
          name:        'authorization-service-password'
          keyVaultUrl: 'https://${keyVaultName}.vault.azure.net/secrets/authorizationServicePassword'
          identity:    'system'
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
          image: '${containerRepositoryName}.azurecr.io/${appName}:default' 
          name:  appName
          env: [
            {
              name: 'RESOURCE_GROUP'
              value: resourceGroup().name
            }
          ]
        }
      ]
    }
  }
}

module keyVaultRoleAssignment 'keyVaultSecretsUserRoleAssignment.bicep' = {
  name:  'keyVaultRoleAssignment' 
  scope: resourceGroup(commonEnvironmentResourceGroup)
  params: {
     keyVaultName: keyVaultName
     prefix:       prefix
     principalId:  containerApp.identity.principalId
  }
}