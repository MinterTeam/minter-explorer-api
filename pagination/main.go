package pagination

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"math"
	"net/http"
	"strconv"
)

type Service struct {
	Pager      *urlvalues.Pager
	Request    *http.Request
	RequestURL string
	Total      int
}

const DefaultLimit = "50"

func NewService(request *http.Request) Service {
	values := urlvalues.Values(request.URL.Query())
	values.SetDefault("limit", DefaultLimit)

	return Service{
		Pager:      values.Pager(),
		Request:    request,
		RequestURL: fmt.Sprintf("%s%s%s", request.URL.Scheme, request.Host, request.URL.Path),
	}
}

func (service Service) ApplyFilter(query *orm.Query) *orm.Query {
	paginatedQuery, err := service.Pager.Pagination(query)
	helpers.CheckErr(err)

	return paginatedQuery
}

func (service Service) GetNextPageLink() *string {
	if service.GetLastPage() == service.GetCurrentPage() {
		return nil
	}

	nextPage := strconv.Itoa(service.GetCurrentPage() + 1)
	query := service.Request.URL.Query()
	query.Set("page", nextPage)

	link := fmt.Sprintf("%s?%s", service.RequestURL, query.Encode())
	return &link
}

func (service Service) GetLastPageLink() *string {
	lastPage := strconv.Itoa(service.GetLastPage())
	query := service.Request.URL.Query()
	query.Set("page", lastPage)

	link := fmt.Sprintf("%s?%s", service.RequestURL, query.Encode())
	return &link
}

func (service Service) GetPrevPageLink() *string {
	if service.GetCurrentPage() == 1 {
		return nil
	}

	prevPage := strconv.Itoa(service.GetCurrentPage() - 1)
	query := service.Request.URL.Query()
	query.Set("page", prevPage)

	link := fmt.Sprintf("%s?%s", service.RequestURL, query.Encode())
	return &link
}

func (service Service) GetFirstPageLink() *string {
	query := service.Request.URL.Query()
	query.Set("page", "1")

	link := fmt.Sprintf("%s?%s", service.RequestURL, query.Encode())
	return &link
}

func (service Service) GetPath() string {
	return service.RequestURL
}

func (service Service) GetLastPage() int {
	return int(math.Ceil(float64(service.Total) / float64(service.Pager.Limit)))
}

func (service Service) GetCurrentPage() int {
	return service.Pager.GetPage()
}

func (service Service) GetPerPage() int {
	return service.Pager.GetLimit()
}
