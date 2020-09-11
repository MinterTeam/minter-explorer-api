package tools

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/go-pg/pg/v9/orm"
	"github.com/go-pg/urlstruct"
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	Pager      *urlstruct.Pager
	Request    *http.Request
	RequestURL string
	Total      int
}

func NewPagination(request *http.Request) Pagination {
	query := request.URL.Query()
	if query.Get("limit") == "" {
		query.Set("limit", strconv.FormatInt(config.DefaultPaginationLimit, 10))
	}

	limit, _ := strconv.ParseUint(query.Get("limit"), 10, 64)
	if limit > config.MaxPaginationLimit {
		query.Set("limit", strconv.FormatInt(config.MaxPaginationLimit, 10))
	}

	pager := urlstruct.NewPager(query)
	pager.MaxOffset = config.MaxPaginationOffset
	pager.MaxLimit = config.MaxPaginationLimit

	return Pagination{
		Pager:      pager,
		Request:    request,
		RequestURL: fmt.Sprintf("https://%s%s", request.Host, request.URL.Path), // TODO: fix request url
	}
}

func (pagination Pagination) Filter(query *orm.Query) (*orm.Query, error) {
	return query.Limit(pagination.Pager.GetLimit()).Offset(pagination.Pager.GetOffset()), nil
}

func (pagination Pagination) GetNextPageLink() *string {
	if pagination.GetLastPage() == 0 || pagination.GetLastPage() == pagination.GetCurrentPage() {
		return nil
	}

	nextPage := strconv.Itoa(pagination.GetCurrentPage() + 1)
	query := pagination.Request.URL.Query()
	query.Set("page", nextPage)

	link := fmt.Sprintf("%s?%s", pagination.RequestURL, query.Encode())
	return &link
}

func (pagination Pagination) GetLastPageLink() *string {
	lastPage := "1"
	if pagination.GetLastPage() != 0 {
		lastPage = strconv.Itoa(pagination.GetLastPage())
	}

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
