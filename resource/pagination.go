package resource

import (
	"github.com/MinterTeam/minter-explorer-api/pagination"
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
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int    `json:"per_page"`
	Total       int    `json:"total"`
}

func TransformPaginatedCollection(collection interface{}, resource Interface, paginationService pagination.Service) PaginationResource {
	result := TransformCollection(collection, resource)

	return PaginationResource{
		Data: result,
		Links: PaginationLinksResource{
			First: paginationService.GetFirstPageLink(),
			Last:  paginationService.GetLastPageLink(),
			Prev:  paginationService.GetPrevPageLink(),
			Next:  paginationService.GetNextPageLink(),
		},
		Meta: PaginationMetaResource{
			CurrentPage: paginationService.GetCurrentPage(),
			LastPage:    paginationService.GetLastPage(),
			Path:        paginationService.GetPath(),
			PerPage:     paginationService.GetPerPage(),
			Total:       paginationService.Total,
		},
	}
}
