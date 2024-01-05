package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/go-logr/logr"
	"net/http"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"golang.org/x/net/context"
)

func ReadItem[T any](
	ctx       context.Context,
	id        string,
	pk        azcosmos.PartitionKey,
	container *azcosmos.ContainerClient,
	opts      *azcosmos.ItemOptions,
) (azcosmos.ItemResponse, T, bool, error) {
	if container == nil {
		panic("container cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	var item T
	var resError *azcore.ResponseError
	res, err := container.ReadItem(ctx, pk, id, opts)
	if errors.As(err, &resError) {
		if resError.StatusCode == http.StatusNotFound {
			return res, item, false, nil
		}
	} else if err != nil {
		return res, item, false, err
	}
	err = json.Unmarshal(res.Value, &item)
	return res, item, true, err
}

func CreateItem[T any](
	ctx       context.Context,
	item      *T,
	pk        azcosmos.PartitionKey,
	container *azcosmos.ContainerClient,
	opts      *azcosmos.ItemOptions,
) (azcosmos.ItemResponse, error) {
	if container == nil {
		panic("container cannot be nil")
	}
	bytes, err := json.Marshal(item)
	if err != nil {
		return azcosmos.ItemResponse{}, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return container.CreateItem(ctx, pk, bytes, opts)
}

func ReplaceItem[T any](
	ctx       context.Context,
	op        func(*T) (bool, error),
	id        string,
	pk        azcosmos.PartitionKey,
	container *azcosmos.ContainerClient,
	opts      *azcosmos.ItemOptions,
) (azcosmos.ItemResponse, bool, error) {
	if container == nil {
		panic("container cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	var resError *azcore.ResponseError
	res, err := container.ReadItem(ctx, pk, id, nil)
	if errors.As(err, &resError) {
		if resError.StatusCode == http.StatusNotFound {
			return res, false, nil
		}
	} else if err != nil {
		return res, false, err
	}
	var item T
	err = json.Unmarshal(res.Value, &item)
	if err != nil {
		return res, true, err
	}
	if changed, err := op(&item); err != nil {
		return res, true, err
	} else if !changed {
		return res, true, nil
	}
	bytes, err := json.Marshal(item)
	if err != nil {
		return res, true, err
	}
	if opts != nil {
		if opts.IfMatchEtag == nil {
			opts.IfMatchEtag = &res.ETag
		}
	} else {
		opts = &azcosmos.ItemOptions{
			IfMatchEtag: &res.ETag,
		}
	}
	res, err = container.ReplaceItem(ctx, pk, id, bytes, opts)
	return res, true, err
}

func DeleteItem(
	ctx       context.Context,
	id        string,
	pk        azcosmos.PartitionKey,
	container *azcosmos.ContainerClient,
	opts      *azcosmos.ItemOptions,
) (azcosmos.ItemResponse, bool, error) {
	if container == nil {
		panic("container cannot be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	var resError *azcore.ResponseError
	res, err := container.DeleteItem(ctx, pk, id, opts)
	if errors.As(err, &resError) {
		if resError.StatusCode == http.StatusNotFound {
			return res, false, nil
		}
	}
	return res, true, err
}

func Container(
	azClient    *azcosmos.Client,
	databaseId  string,
	containerId string,
) *azcosmos.ContainerClient {
	if azClient == nil {
		panic("azClient cannot be nil")
	}
	database, err := azClient.NewDatabase(databaseId)
	if err != nil {
		panic(fmt.Errorf("error creating %q db client: %w", databaseId, err))
	}
	container, err := database.NewContainer(containerId)
	if err != nil {
		panic(fmt.Errorf("error creating %q container client: %w", containerId, err))
	}
	return container
}

func ProvisionDatabase(
	ctx      context.Context,
	azClient *azcosmos.Client,
	cfg      *Config,
	logger   logr.Logger,
) error {
	if azClient == nil {
		panic("azClient cannot be nil")
	}
	if cfg == nil {
		panic("cfg cannot be nil")
	}
	var options *azcosmos.CreateDatabaseOptions
	th := cfg.ThroughputChoice()
	if th != nil {
		options = &azcosmos.CreateDatabaseOptions{
			ThroughputProperties: th,
		}
	}

	// Create database
	var resError *azcore.ResponseError
	_, err := azClient.CreateDatabase(ctx, azcosmos.DatabaseProperties{
		ID: cfg.Name,
	}, options)
	if errors.As(err, &resError) {
		if resError.StatusCode != http.StatusConflict {
			return err
		}
		logger.Info("database exists", "name", cfg.Name)
	} else if err != nil {
		return err
	} else {
		logger.Info("database created", "name", cfg.Name)
	}
	dbClient, err := azClient.NewDatabase(cfg.Name)
	if err != nil {
		return nil
	}

	// Update database
	if th != nil {
		res, err := dbClient.ReadThroughput(ctx, nil)
		if err == nil {
			tp := res.ThroughputProperties
			if !reflect.DeepEqual(th, tp) {
				_, err = dbClient.ReplaceThroughput(ctx, *th, nil)
			}
		}
		if err != nil {
			return err
		}
	}

	// Provision containers
	for _, cnt := range cfg.Containers {
		if err = ProvisionContainer(ctx, dbClient, &cnt, logger); err != nil {
			return err
		}
	}
	return nil
}

func ProvisionContainer(
	ctx      context.Context,
	database *azcosmos.DatabaseClient,
	cfg      *ContainerConfig,
	logger   logr.Logger,
) error {
	if database == nil {
		panic("database cannot be nil")
	}
	if cfg == nil {
		panic("cfg cannot be nil")
	}
	var options *azcosmos.CreateContainerOptions
	th := cfg.ThroughputChoice()
	if th != nil {
		options = &azcosmos.CreateContainerOptions{
			ThroughputProperties: th,
		}
	}

	// Create container
	var resError *azcore.ResponseError
	uk := cfg.UniqueKeys
	_, err := database.CreateContainer(ctx, azcosmos.ContainerProperties{
		ID:                     cfg.Name,
		PartitionKeyDefinition: cfg.PartitionKey,
		IndexingPolicy:         cfg.Indexes,
		UniqueKeyPolicy:        uk,
	}, options)
	if errors.As(err, &resError) {
		if resError.StatusCode != http.StatusConflict {
			return err
		}
		logger.Info("container exists", "database", database.ID(), "name", cfg.Name)
	} else if err != nil {
		return err
	} else {
		logger.Info("container created", "database", database.ID(), "name", cfg.Name)
		return nil
	}

	// Update container
	if uk == nil && th == nil {
		return nil
	}
	container, err := database.NewContainer(cfg.Name)
	if err != nil {
		return err
	}
	res, err := container.Read(ctx, nil)
	if err != nil {
		return err
	}
	if uk != nil {
		cp := res.ContainerProperties
		if !reflect.DeepEqual(uk, cp.UniqueKeyPolicy) {
			cp.UniqueKeyPolicy = uk
			_, err = container.Replace(ctx, *cp, nil)
			if err != nil {
				return err
			}
		}
	}
	if th != nil {
		res, err := container.ReadThroughput(ctx, nil)
		if err == nil {
			tp := res.ThroughputProperties
			if !reflect.DeepEqual(th, tp) {
				_, err = container.ReplaceThroughput(ctx, *th, nil)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}