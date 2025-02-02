package heimdall

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ledgerwatch/log/v3"
)

type entityFetcher[TEntity Entity] interface {
	FetchLastEntityId(ctx context.Context) (uint64, error)
	FetchEntitiesRange(ctx context.Context, idRange ClosedRange) ([]TEntity, error)
}

type entityFetcherImpl[TEntity Entity] struct {
	name string

	fetchLastEntityId func(ctx context.Context) (int64, error)
	fetchEntity       func(ctx context.Context, id int64) (TEntity, error)
	fetchEntitiesPage func(ctx context.Context, page uint64, limit uint64) ([]TEntity, error)

	logger log.Logger
}

func newEntityFetcher[TEntity Entity](
	name string,
	fetchLastEntityId func(ctx context.Context) (int64, error),
	fetchEntity func(ctx context.Context, id int64) (TEntity, error),
	fetchEntitiesPage func(ctx context.Context, page uint64, limit uint64) ([]TEntity, error),
	logger log.Logger,
) entityFetcher[TEntity] {
	return &entityFetcherImpl[TEntity]{
		name:              name,
		fetchLastEntityId: fetchLastEntityId,
		fetchEntity:       fetchEntity,
		fetchEntitiesPage: fetchEntitiesPage,
		logger:            logger,
	}
}

func (f *entityFetcherImpl[TEntity]) FetchLastEntityId(ctx context.Context) (uint64, error) {
	id, err := f.fetchLastEntityId(ctx)
	return uint64(id), err
}

func (f *entityFetcherImpl[TEntity]) FetchEntitiesRange(ctx context.Context, idRange ClosedRange) ([]TEntity, error) {
	count := idRange.Len()

	const batchFetchThreshold = 100
	if (count > batchFetchThreshold) && (f.fetchEntitiesPage != nil) {
		allEntities, err := f.FetchAllEntities(ctx)
		if err != nil {
			return nil, err
		}
		startIndex := idRange.Start - 1
		return allEntities[startIndex : startIndex+count], nil
	}

	return f.FetchEntitiesRangeSequentially(ctx, idRange)
}

func (f *entityFetcherImpl[TEntity]) FetchEntitiesRangeSequentially(ctx context.Context, idRange ClosedRange) ([]TEntity, error) {
	return ClosedRangeMap(idRange, func(id uint64) (TEntity, error) {
		return f.fetchEntity(ctx, int64(id))
	})
}

func (f *entityFetcherImpl[TEntity]) FetchAllEntities(ctx context.Context) ([]TEntity, error) {
	// TODO: once heimdall API is fixed to return sorted items in pages we can only fetch
	//
	//	the new pages after lastStoredCheckpointId using the checkpoints/list paging API
	//	(for now we have to fetch all of them)
	//	and also remove sorting we do after fetching

	var entities []TEntity

	fetchStartTime := time.Now()
	progressLogTicker := time.NewTicker(30 * time.Second)
	defer progressLogTicker.Stop()

	for page := uint64(1); ; page++ {
		entitiesPage, err := f.fetchEntitiesPage(ctx, page, 10_000)
		if err != nil {
			return nil, err
		}
		if len(entitiesPage) == 0 {
			break
		}

		for _, entity := range entitiesPage {
			entities = append(entities, entity)
		}

		select {
		case <-progressLogTicker.C:
			f.logger.Debug(
				heimdallLogPrefix(fmt.Sprintf("%s progress", f.name)),
				"page", page,
				"len", len(entities),
			)
		default:
			// carry-on
		}
	}

	slices.SortFunc(entities, func(e1, e2 TEntity) int {
		n1 := e1.BlockNumRange().Start
		n2 := e2.BlockNumRange().Start
		return cmp.Compare(n1, n2)
	})

	f.logger.Debug(
		heimdallLogPrefix(fmt.Sprintf("%s done", f.name)),
		"len", len(entities),
		"duration", time.Since(fetchStartTime),
	)

	return entities, nil
}
