package db

import (
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type (
	Config struct {
		Name          string
		ConnectionUri string
		Timeout       time.Duration
		Containers    []ContainerConfig
		Throughput    *ThroughputConfig
	}

	ThroughputConfig struct {
		Manual *struct{
			Throughput int32
		}
		AutoScale *struct{
			StartingMaxThroughput int32
		}
		AutoScaleIncrement *struct{
			StartingMaxThroughput int32
			IncrementPercentage   int32
		}
	}

	ContainerConfig struct {
		Name         string
		PartitionKey azcosmos.PartitionKeyDefinition
		UniqueKeys   *azcosmos.UniqueKeyPolicy
		Indexes      *azcosmos.IndexingPolicy
		Conflicts    *azcosmos.ConflictResolutionPolicy
		Throughput   *ThroughputConfig
	}
)

func (c *Config) Validate() error {
	if th := c.Throughput; th != nil {
		cnt := 0
		if th.Manual != nil {
			cnt++
		}
		if th.AutoScale != nil {
			cnt++
		}
		if th.AutoScaleIncrement != nil {
			cnt++
		}
		switch {
		case cnt > 1:
			return errors.New("only one throughput choice is allowed")
		case cnt == 0:
			return errors.New("throughput choice is missing")
		}
	}
	return nil
}

func (c *Config) ThroughputChoice() *azcosmos.ThroughputProperties {
	return toThroughput(c.Throughput)
}

func (c *ContainerConfig) ThroughputChoice() *azcosmos.ThroughputProperties {
	return toThroughput(c.Throughput)
}


func toThroughput(th *ThroughputConfig) *azcosmos.ThroughputProperties {
	if th != nil {
		if manual := th.Manual; manual != nil {
			p := azcosmos.NewManualThroughputProperties(manual.Throughput)
			return &p
		}
		if autoScale := th.AutoScale; autoScale != nil {
			p := azcosmos.NewAutoscaleThroughputProperties(autoScale.StartingMaxThroughput)
			return &p
		}
		if autoScaleInc := th.AutoScaleIncrement; autoScaleInc != nil {
			p := azcosmos.NewAutoscaleThroughputPropertiesWithIncrement(
				autoScaleInc.StartingMaxThroughput, autoScaleInc.IncrementPercentage)
			return &p
		}
	}
	return nil
}