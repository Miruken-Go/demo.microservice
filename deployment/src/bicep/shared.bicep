param containerRepositoryName string
param location            string 

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' = {
  name:     containerRepositoryName
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    adminUserEnabled: true
  }
}
