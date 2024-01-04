package db

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type (
	Config struct {
		Name          string
		ConnectionUri string
		Timeout       time.Duration
		Containers    []ContainerConfig
	}

	ContainerConfig struct {
		Name         string
		PartitionKey azcosmos.PartitionKeyDefinition
		UniqueKeys   *azcosmos.UniqueKeyPolicy
		Indexes      *azcosmos.IndexingPolicy
		Conflicts	 *azcosmos.ConflictResolutionPolicy
	}
)

