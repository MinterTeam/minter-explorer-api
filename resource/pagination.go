package resource

import (
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"sync"
)

type PaginationResource struct {
	Data  []Interface             `json:"data"`
	Links PaginationLinksResource `json:"links"`
	Meta  PaginationMetaResource  `json:"meta"`
}

type PaginationLinksResource struct {
	First *string `json:"first"`
	Last  *string `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

type PaginationMetaResource struct {
	CurrentPage int                    `json:"current_page"`
	LastPage    int                    `json:"last_page"`
	Path        string                 `json:"path"`
	PerPage     int                    `json:"per_page"`
	Total       int                    `json:"total"`
	Additional  map[string]interface{} `json:"additional,omitempty"`
}

func TransformPaginatedCollection(collection interface{}, resource Interface, pagination tools.Pagination) PaginationResource {
	return transformPaginatedCollection(collection, resource, pagination, nil)
}

func TransformPaginatedCollectionWithCallback(collection interface{}, resource Interface, pagination tools.Pagination, callbackFunc func(model ParamInterface) ParamsInterface) PaginationResource {
	models := makeItemsFromModelsCollection(collection)
	result := make([]Interface, len(models))

	wg := new(sync.WaitGroup)
	for i := range models {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			defer errors.Recovery()
			result[i] = resource.Transform(models[i], callbackFunc(models[i])...)
		}(i, wg)
	}
	wg.Wait()

	return PaginationResource{
		Data: result,
		Links: PaginationLinksResource{
			First: pagination.GetFirstPageLink(),
			Last:  pagination.GetLastPageLink(),
			Prev:  pagination.GetPrevPageLink(),
			Next:  pagination.GetNextPageLink(),
		},
		Meta: PaginationMetaResource{
			CurrentPage: pagination.GetCurrentPage(),
			LastPage:    pagination.GetLastPage(),
			Path:        pagination.GetPath(),
			PerPage:     pagination.GetPerPage(),
			Total:       pagination.Total,
		},
	}
}

func transformPaginatedCollection(collection interface{}, resource Interface, pagination tools.Pagination, additional map[string]interface{}) PaginationResource {
	result := TransformCollection(collection, resource)

	return PaginationResource{
		Data: result,
		Links: PaginationLinksResource{
			First: pagination.GetFirstPageLink(),
			Last:  pagination.GetLastPageLink(),
			Prev:  pagination.GetPrevPageLink(),
			Next:  pagination.GetNextPageLink(),
		},
		Meta: PaginationMetaResource{
			CurrentPage: pagination.GetCurrentPage(),
			LastPage:    pagination.GetLastPage(),
			Path:        pagination.GetPath(),
			PerPage:     pagination.GetPerPage(),
			Total:       pagination.Total,
			Additional:  additional,
		},
	}
}
