param prefix       string
param keyVaultName string
param principalId  string

resource keyVault 'Microsoft.KeyVault/vaults@2023-02-01' existing = {
  name:  keyVaultName
}

resource containerApp_keyvault_role 'Microsoft.Authorization/roleAssignments@2020-08-01-preview' = {
  name:  guid('containerApp_keyvault_role', prefix)
  scope: keyVault
  properties: {
    roleDefinitionId: '/providers/Microsoft.Authorization/roleDefinitions/4633458b-17de-408a-b874-0445c86b69e6'
    description:      'Key Vault Secrets User'
    principalType:    'ServicePrincipal'
    principalId:      principalId
  }
}
