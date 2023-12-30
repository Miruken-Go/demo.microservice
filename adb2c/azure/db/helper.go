package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"golang.org/x/net/context"
)

func Container(
	azure       *azcosmos.Client,
	databaseId  string,
	containerId string,
) *azcosmos.ContainerClient {
	database, err := azure.NewDatabase(databaseId)
	if err != nil {
		panic(fmt.Errorf("error creating %q db client: %w", databaseId, err))
	}
	container, err := database.NewContainer(containerId)
	if err != nil {
		panic(fmt.Errorf("error creating %q container client: %w", containerId, err))
	}
	return container
}

func ReadItem[T any](
	ctx       context.Context,
	id        string,
	pk        azcosmos.PartitionKey,
	container *azcosmos.ContainerClient,
	opts      *azcosmos.ItemOptions,
) (azcosmos.ItemResponse, T, bool, error) {
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