package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-logr/logr"
	"github.com/miruken-go/miruken/promise"
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
) *promise.Promise[struct{}] {
	if azClient == nil {
		panic("azClient cannot be nil")
	}
	if cfg == nil {
		panic("cfg cannot be nil")
	}

	return promise.New(nil, func(
		resolve func(struct{}), reject func(error), onCancel func(func())) {

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
				reject(err)
				return
			}
			logger.Info("database exists", "database", cfg.Name)
		} else if err != nil {
			reject(err)
			return
		} else {
			logger.Info("database created", "database", cfg.Name)
		}
		dbClient, err := azClient.NewDatabase(cfg.Name)
		if err != nil {
			reject(err)
			return
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
				reject(err)
				return
			}
		}

		// Provision containers
		containers := cfg.Containers
		if len(containers) == 0 {
			resolve(struct{}{})
		}
		promises := make([]*promise.Promise[struct{}], len(containers))
		for i := range containers {
			promises[i] = ProvisionContainer(ctx, dbClient, &containers[i], logger)
		}
		if _, err := promise.All(ctx, promises...).Await(); err != nil {
			reject(err)
			return
		}
		resolve(struct{}{})
	})
}

func ProvisionContainer(
	ctx      context.Context,
	database *azcosmos.DatabaseClient,
	cfg      *ContainerConfig,
	logger   logr.Logger,
) *promise.Promise[struct{}] {
	if database == nil {
		panic("database cannot be nil")
	}
	if cfg == nil {
		panic("cfg cannot be nil")
	}

	return promise.New(nil, func(
		resolve func(struct{}), reject func(error), onCancel func(func())) {
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
				reject(err)
				return
			}
			logger.Info("container exists",
				"database", database.ID(), "container", cfg.Name)
		} else if err != nil {
			reject(err)
			return
		} else {
			logger.Info("container created",
				"database", database.ID(), "container", cfg.Name)
			resolve(struct{}{})
			return
		}

		// Update container
		if uk == nil && th == nil {
			resolve(struct{}{})
			return
		}
		container, err := database.NewContainer(cfg.Name)
		if err != nil {
			reject(err)
			return
		}
		res, err := container.Read(ctx, nil)
		if err != nil {
			reject(err)
			return
		}
		if uk != nil {
			cp := res.ContainerProperties
			if !reflect.DeepEqual(uk, cp.UniqueKeyPolicy) {
				cp.UniqueKeyPolicy = uk
				_, err = container.Replace(ctx, *cp, nil)
				if err != nil {
					reject(err)
					return
				}
				logger.Info("container indexes updated",
					"database", database.ID(), "container", cfg.Name)
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
				reject(err)
				return
			}
			logger.Info("container throughput updated",
				"database", database.ID(), "container", cfg.Name)
		}
		resolve(struct{}{})
	})
}