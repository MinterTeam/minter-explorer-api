package tools

import (
"fmt"
"github.com/MinterTeam/minter-explorer-api/helpers"
"github.com/go-pg/pg/orm"
"github.com/go-pg/pg/urlvalues"
"math"
"net/http"
"strconv"
)

type Pagination struct {
	Pager      *urlvalues.Pager
	Request    *http.Request
	RequestURL string
	Total      int
}

const DefaultLimit = 50

func NewPagination(request *http.Request) Pagination {
	values := urlvalues.Values(request.URL.Query())
	values.SetDefault("limit", strconv.Itoa(DefaultLimit))

	return Pagination{
		Pager:      values.Pager(),
		Request:    request,
		RequestURL: fmt.Sprintf("%s%s%s", request.URL.Scheme, request.Host, request.URL.Path),
	}
}

func (pagination Pagination) ApplyFilter(query *orm.Query) *orm.Query {
	paginatedQuery, err := pagination.Pager.Pagination(query)
	helpers.CheckErr(err)

	return paginatedQuery
}

func (pagination Pagination) GetNextPageLink() *string {
	if pagination.GetLastPage() == pagination.GetCurrentPage() {
		return nil
	}

	nextPage := strconv.Itoa(pagination.GetCurrentPage() + 1)
	query := pagination.Request.URL.Query()
	query.Set("page", nextPage)

	link := fmt.Sprintf("%s?%s", pagination.RequestURL, query.Encode())
	return &link
}

func (pagination Pagination) GetLastPageLink() *string {
	lastPage := strconv.Itoa(pagination.GetLastPage())
	query := pagination.Request.URL.Query()
	query.Set("page", lastPage)

	link := fmt.Sprintf("%s?%s", pagination.RequestURL, query.Encode())
	return &link
}

func (pagination Pagination) GetPrevPageLink() *string {
	if pagination.GetCurrentPage() == 1 {
		return nil
	}

	prevPage := strconv.Itoa(pagination.GetCurrentPage() - 1)
	query := pagination.Request.URL.Query()
	query.Set("page", prevPage)

	link := fmt.Sprintf("%s?%s", pagination.RequestURL, query.Encode())
	return &link
}

func (pagination Pagination) GetFirstPageLink() *string {
	query := pagination.Request.URL.Query()
	query.Set("page", "1")

	link := fmt.Sprintf("%s?%s", pagination.RequestURL, query.Encode())
	return &link
}

func (pagination Pagination) GetPath() string {
	return pagination.RequestURL
}

func (pagination Pagination) GetLastPage() int {
	return int(math.Ceil(float64(pagination.Total) / float64(pagination.Pager.Limit)))
}

func (pagination Pagination) GetCurrentPage() int {
	return pagination.Pager.GetPage()
}

func (pagination Pagination) GetPerPage() int {
	return pagination.Pager.GetLimit()
}
