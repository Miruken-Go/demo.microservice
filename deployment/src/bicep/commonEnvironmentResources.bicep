param keyVaultName string
param location     string 

/////////////////////////////////////////////////////////////////////////////////////
// KeyVault
/////////////////////////////////////////////////////////////////////////////////////

resource keyVault 'Microsoft.KeyVault/vaults@2023-02-01' = {
  name: keyVaultName
  location: location
  properties: {
    enabledForDeployment:         true
    enabledForTemplateDeployment: true
    enabledForDiskEncryption:     true
    enableRbacAuthorization:      true
    tenantId: subscription().tenantId
    sku: {
      name:   'standard'
      family: 'A'
    }
    networkAcls: {
      defaultAction: 'Allow'
      bypass:        'AzureServices'
    }
  }
}
