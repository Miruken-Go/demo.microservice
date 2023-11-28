@secure()
param containerRepositoryPassword string 

param prefix                      string
param appName                     string 
param containerAppName            string
param location                    string 
param containerAppsEnvironmentId  string
param containerRepositoryName     string 
param keyVaultName                string
param keyVaultResourceGroup       string
param secrets                     array

var keyVaultSecrets = [for secret in secrets: {
  name:        secret
  keyVaultUrl: 'https://${keyVaultName}.vault.azure.net/secrets/${secret}'
  identity:    'system'
}]

resource containerApp 'Microsoft.App/containerApps@2023-05-01' ={
  name:     containerAppName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties:{
    managedEnvironmentId: containerAppsEnvironmentId
    configuration: {
      activeRevisionsMode: 'Multiple'
      ingress: {
        targetPort: 8080
        external: true
      }
      secrets: concat([{
          name: 'acr-password'
          value: containerRepositoryPassword
        }], 
        keyVaultSecrets
      )
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
          image: '${containerRepositoryName}.azurecr.io/${appName}:default' 
          name:  appName
          env: [
            {
              name: 'RESOURCE_GROUP'
              value: resourceGroup().name
            }
            {
              name: 'APPLICATION_NAME'
              value: appName
            }
          ]
        }
      ]
    }
  }
}

module keyVaultRoleAssignment 'keyVaultSecretsUserRoleAssignment.bicep' = {
  name:  'keyVaultRoleAssignment' 
  scope: resourceGroup(keyVaultResourceGroup)
  params: {
     keyVaultName:     keyVaultName
     prefix:           prefix
     containerAppName: containerAppName
     principalId:      containerApp.identity.principalId
  }
}

output containerApp object = containerApp
