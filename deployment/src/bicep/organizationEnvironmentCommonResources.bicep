param prefix       string
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

/////////////////////////////////////////////////////////////////////////////////////
// CosmosDb
/////////////////////////////////////////////////////////////////////////////////////

resource name_resource 'Microsoft.DocumentDb/databaseAccounts@2023-11-15-preview' = {
  name:     prefix
  location: location
  kind:     'GlobalDocumentDB'
  properties: {
    databaseAccountOfferType: 'Standard'
    locations: [
      {
        id: '${prefix}-${location}'
        failoverPriority: 0
        locationName: location
      }
    ]
    backupPolicy: {
      type: 'Continuous'
      continuousModeProperties: {
        tier: 'Continuous7Days'
      }
    }
    isVirtualNetworkFilterEnabled: false
    virtualNetworkRules: []
    ipRules: []
    dependsOn: []
    minimalTlsVersion: 'Tls12'
    capabilities: [
      {
        name: 'EnableServerless'
      }
    ]
    enableFreeTier: false
    capacity: {
      totalThroughputLimit: 4000
    }
  }
}

