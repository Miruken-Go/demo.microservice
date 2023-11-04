param prefix                  string
param location                string 

/////////////////////////////////////////////////////////////////////////////////////
// KeyVault
/////////////////////////////////////////////////////////////////////////////////////

resource keyVault 'Microsoft.KeyVault/vaults@2019-09-01' = {
  name: '${prefix}-keyvault'
  location: location
  properties: {
    enabledForDeployment:         true
    enabledForTemplateDeployment: true
    enabledForDiskEncryption:     true
    enableRbacAuthorization:      true
    tenantId: tenantId
    accessPolicies: []
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
// Cosmos
/////////////////////////////////////////////////////////////////////////////////////

@description('Standard envs are dev, qa, prod')
param primaryRegion             string
param secondaryRegion           string
param vespyEventStoreThroughput int
param vespyCentralThroughput    int
param vespyIoTThroughput        int

var locations = useProdSetting ? [
  {
    locationName: primaryRegion
    failoverPriority: 0
    isZoneRedundant: false
  }
  {
    locationName: secondaryRegion
    failoverPriority: 1
    isZoneRedundant: false
  }
] : [
  {
    locationName: primaryRegion
    failoverPriority: 0
    isZoneRedundant: false
  }
]

var virtualNetworkRules = useDevSetting ? [] : [
  {
    id:                               webSiteSubnetId
    ignoreMissingVNetServiceEndpoint: false
  }
  {
    id:                               defaultSubnetId 
    ignoreMissingVNetServiceEndpoint: false
  }
]

var ipRules = useDevSetting ? [] : [
  {
    ipAddressOrRange: '104.42.195.92'
  }
  {
    ipAddressOrRange: '40.76.54.131'
  }
  {
    ipAddressOrRange: '52.176.6.30'
  }
  {
    ipAddressOrRange: '52.169.50.45'
  }
  {
    ipAddressOrRange: '52.187.184.26'
  }
]

resource databaseAccount 'Microsoft.DocumentDB/databaseAccounts@2021-07-01-preview' = {
  name:     '${prefix}-cosmos-account'
  location: location
  kind:     'MongoDB'
  tags: {}
  properties: {
    enableAutomaticFailover:      true
    enableMultipleWriteLocations: false
    databaseAccountOfferType:     'Standard'
    apiProperties: {
      serverVersion: '4.2' 
    }
    consistencyPolicy: {
      defaultConsistencyLevel: 'Strong'
      maxIntervalInSeconds:    5
      maxStalenessPrefix:      100
    }
    locations: locations
    capabilities: [
      {
        name: 'EnableMongo'
      }
      {
        name: 'DisableRateLimitingResponses'
      }
    ]
    virtualNetworkRules:           virtualNetworkRules
    isVirtualNetworkFilterEnabled: useDevSetting ? false : true
    ipRules:                       ipRules
    backupPolicy: {
      type: 'Continuous'
    }
    diagnosticLogSettings: {
      enableFullTextQuery: 'True'
    }
  }
}

resource VespyEventStore 'Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-06-15' = {
  parent: databaseAccount
  name: 'VespyEventStore'
  properties: {
    resource: {
      id: 'VespyEventStore'
    }
    options: {
      autoscaleSettings: {
        maxThroughput: vespyEventStoreThroughput
      }
    }
  }

  resource Event 'collections' = {
    name: 'Event'
    properties: {
      resource: {
        id: 'Event'
        shardKey: {
          EntityId: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                  '_id'
              ]
            }
          }
          {
            key: {
                keys: [
                  'EntityId'
                  'SeqNum'
                ]
            }
            options: {
                unique: true
            }
          }
          {
            key: {
              keys: [
                  'Timestamp'
              ]
            }
          }
          {
            key: {
              keys: [
                  'TxId'
                  'Accepted'
              ]
            }
          }
        ]
      }
    }
  }
}

resource vespyCentral 'Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-06-15' = {
  parent: databaseAccount
  name: 'VespyCentral'
  properties: {
    resource: {
      id: 'VespyCentral'
    }
    options: {
      autoscaleSettings: {
        maxThroughput: vespyCentralThroughput
      }
    }
  }

  resource vespyCentral_customer 'collections' = {
    name: 'Customer'
    properties: {
      resource: {
        id: 'Customer'
        shardKey: {
          _id: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                  '_id'
              ]
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }

  resource vespyCentral_Installation 'collections' = {
    name: 'Installation'
    properties: {
      resource: {
        id: 'Installation'
        shardKey: {
          CustomerId: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                  '_id'
              ]
            }
          }
          {
            key: {
                keys: [
                  'CustomerId'
                  'Name'
                ]
            }
            options: {
                unique: true
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }

  resource vespyCentral_User 'collections' = {
    name: 'User'
    properties: {
      resource: {
        id: 'User'
        shardKey: {
          Email: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                  '_id'
              ]
            }
          }
          {
            key: {
                keys: [
                  'Email'
                ]
            }
            options: {
                unique: true
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }

  resource vespyCentral_Zone 'collections' = {
    name: 'Zone'
    properties: {
      resource: {
        id: 'Zone'
        shardKey: {
          InstallationId: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                  '_id'
              ]
            }
          }
          {
            key: {
                keys: [
                  'InstallationId'
                  'Name'
                ]
            }
            options: {
                unique: true
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }

  resource vespyCentral_Asset 'collections' = {
    name: 'Asset'
    properties: {
      resource: {
        id: 'Asset'
        shardKey: {
          CustomerId: 'Hash'
        } 
        indexes: [
          {
            key: {
              keys: [
                '_id'
              ]
            }
          }
          {
            key:{
              keys: [
                'CustomerId'
                'Epc'
              ]
            }
            options: {
              unique: true
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }
}

resource vespyIoT 'Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-06-15' = {
  parent: databaseAccount
  name: 'VespyIoT'
  properties: {
    resource: {
      id: 'VespyIoT'
    }
    options: useProdSetting ? {
      autoscaleSettings: {
        maxThroughput: vespyIoTThroughput
      }
    } : {
      throughput: vespyIoTThroughput
    }
  }

  resource vespyIoT_Device 'collections' = {
    name: 'Device'
    properties: {
      resource: {
        id: 'Device'
        shardKey: {
          Type: 'Hash'
        }
        indexes: [
          {
            key: {
              keys: [
                '_id'
              ]
            }
          }
          {
            key:{
              keys: [
                'Type'
                'SerialNumber'
              ]
            }
            options: {
              unique: true
            }
          }
          {
            key: {
              keys: [
                  '$**'
              ]
            }
          }
        ]
      }
    }
  }
}

resource vennachassisConnectionString 'Microsoft.KeyVault/vaults/secrets@2019-09-01' = {
  parent: keyVault
  name: 'vennachassisConnectionString'
  properties: {
    value: databaseAccount.listConnectionStrings().connectionStrings[0].connectionString
  }
}

//I want this metric in all environments
resource cosmos_ru_alert 'microsoft.insights/metricalerts@2018-03-01' = {
  name:     '${prefix}-cosmos-ru-alert'
  location: 'global'
  properties: {
    description: 'Alert when cosmos ru\'s reach 80%'
    severity:    0
    enabled:     true
    scopes: [
      databaseAccount.id
    ]
    evaluationFrequency: 'PT1M'
    windowSize:          'PT5M'
    criteria: {
      allOf: [
        {
          threshold:       80
          name:            'Metric1'
          metricNamespace: 'Microsoft.DocumentDB/databaseAccounts'
          metricName:      'NormalizedRUConsumption'
          operator:        'GreaterThan'
          timeAggregation: 'Maximum'
          criterionType:   'StaticThresholdCriterion'
        }
      ]
      'odata.type': 'Microsoft.Azure.Monitor.SingleResourceMultipleMetricCriteria'
    }
    autoMitigate:         true
    targetResourceType:   'Microsoft.DocumentDB/databaseAccounts'
    targetResourceRegion: primaryRegion
    actions: [
      {
        actionGroupId: on_call_actionGroup.id
      }
    ]
  }
}
