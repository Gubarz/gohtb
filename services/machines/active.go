package machines

import (
	"context"
	"errors"
	"strings"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
)

const (
	SortByReleaseDate = v4Client.GetMachinePaginatedParamsSortBy("release-date")
	SortByName        = v4Client.GetMachinePaginatedParamsSortBy("name")
	SortByUserOwns    = v4Client.GetMachinePaginatedParamsSortBy("user-owns")
	SortBySystemOwns  = v4Client.GetMachinePaginatedParamsSortBy("system-owns")
	SortByDifficulty  = v4Client.GetMachinePaginatedParamsSortBy("user-difficulty")
	SortByRating      = v4Client.GetMachinePaginatedParamsSortBy("rating")
)

type ActiveQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        v4Client.GetMachinePaginatedParamsSortBy
	sortType      v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       v4Client.Keyword
}

// ListActive creates a new query for active machines.
// This returns an ActiveQuery that can be chained with filtering and pagination methods.
// Active machines are machines that are currently available.
//
// Example:
//
//	query := client.Machines.ListActive()
//	machines, err := query.ByDifficulty("Hard").ByOS("Linux").Results(ctx)
func (s *Service) ListActive() *ActiveQuery {
	return &ActiveQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

// Next moves to the next page in the pagination sequence.
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	nextPage := query.Next().Results(ctx)
func (q *ActiveQuery) Next() *ActiveQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	prevPage := query.Previous().Results(ctx)
func (q *ActiveQuery) Previous() *ActiveQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// Page sets the specific page number for pagination.
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.Page(3).Results(ctx)
func (q *ActiveQuery) Page(n int) *ActiveQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.PerPage(50).Results(ctx)
func (q *ActiveQuery) PerPage(n int) *ActiveQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// ByCompleted filters machines by completion status.
// Valid values are "Completed" and "InComplete".
// Returns a new ActiveQuery that can be further chained.
// By default, it does not filter by completion status and
// includes both completed and incomplete machines.
//
// Example:
//
//	completed := query.ByCompleted("Completed").Results(ctx)
//	incomplete := query.ByCompleted("InComplete").Results(ctx)
func (q *ActiveQuery) ByCompleted(val string) *ActiveQuery {
	qc := ptr.Clone(q)
	qc.showCompleted = strings.ToLower(val)
	return qc
}

// ByOS filters machines by operating system.
// Valid values include "Linux" and "Windows".
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//		linuxMachines := query.ByOS("Linux").Results(ctx)
//		windowsMachines := query.ByOS("Windows").Results(ctx)
//	 linuxAndWindowsMachines := query.ByOS("Linux").ByOS("Windows").Results(ctx)
func (q *ActiveQuery) ByOS(val string) *ActiveQuery {
	return q.ByOSList(val)
}

// ByOSList filters machines by multiple operating systems.
// Valid values include "Linux" and "Windows".
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.ByOSList("Linux", "Windows").Results(ctx)
func (q *ActiveQuery) ByOSList(val ...string) *ActiveQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.os = append(append([]string{}, q.os...), lowercased...)
	return qc
}

// ByDifficultyList filters machines by multiple difficulty levels.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.ByDifficultyList("Hard", "Insane").Results(ctx)
func (q *ActiveQuery) ByDifficultyList(val ...string) *ActiveQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.difficulty = append(append([]string{}, q.difficulty...), lowercased...)
	return qc
}

// ByDifficulty filters machines by difficulty level.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//		hardMachines := query.ByDifficulty("Hard").Results(ctx)
//		easyMachines := query.ByDifficulty("Easy").Results(ctx)
//	 meduimAndInsaneMachines := query.ByDifficulty("Medium").ByDifficulty("Insane").Results(ctx)
func (q *ActiveQuery) ByDifficulty(val string) *ActiveQuery {
	return q.ByDifficultyList(val)
}

// SortedBy sets the field to sort results by.
// Valid values include "release-date", "name", "user-owns", "system-owns", "rating", "user-difficulty".
// Returns a new UnreleasedQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	machines := query.SortedBy("name").Ascending().Results(ctx)
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *ActiveQuery) SortedBy(field string) *ActiveQuery {
	qc := ptr.Clone(q)
	sortBy := v4Client.GetMachinePaginatedParamsSortBy(strings.ToLower(field))
	qc.sortBy = sortBy
	return qc
}

func (q *ActiveQuery) sort(val v4Client.GetMachinePaginatedParamsSortBy, order v4Client.GetMachinePaginatedParamsSortType) *ActiveQuery {
	qc := ptr.Clone(q)
	qc.sortBy = val
	qc.sortType = order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Ascending().Results(ctx)
func (q *ActiveQuery) Ascending() *ActiveQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetMachinePaginatedParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new ActiveQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *ActiveQuery) Descending() *ActiveQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetMachinePaginatedParamsSortType("desc"))
}

func (q *ActiveQuery) fetchResults(ctx context.Context) (MachinePaginatedResponse, error) {
	params := &v4Client.GetMachinePaginatedParams{
		PerPage:  &q.perPage,
		Page:     &q.page,
		SortBy:   &q.sortBy,
		SortType: &q.sortType,
		Keyword:  &q.keyword,
	}

	if len(q.difficulty) > 0 {
		d := q.difficulty
		params.Difficulty = &d
	}

	if len(q.os) > 0 {
		o := q.os
		params.Os = &o
	}

	if q.showCompleted != "" {
		sc := v4Client.GetMachinePaginatedParamsShowCompleted(q.showCompleted)
		params.ShowCompleted = &sc
	}

	resp, err := q.client.V4().GetMachinePaginated(q.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return MachinePaginatedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachinePaginatedResponse)
	if err != nil {
		return MachinePaginatedResponse{ResponseMeta: meta}, err
	}

	return MachinePaginatedResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// Results executes the query and returns the current page of unreleased machines.
// This method should be called last in the query chain to fetch the actual data.
//
// Example:
//
//	machines, err := client.Machines.ListActive().
//		ByDifficulty("Hard").
//		ByOS("Linux").
//		Page(1).
//		Results(ctx)
func (q *ActiveQuery) Results(ctx context.Context) (MachinePaginatedResponse, error) {
	return q.fetchResults(ctx)
}

// AllResults executes the query and returns all pages of unreleased machines.
// This method automatically paginates through all available results.
// Use with caution for large datasets as it may consume significant memory.
//
// Example:
//
//	allMachines, err := client.Machines.ListActive().
//		ByDifficulty("Hard").
//		AllResults(ctx)
func (q *ActiveQuery) AllResults(ctx context.Context) (MachinePaginatedResponse, error) {
	var all []MachineData
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return resp, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return MachinePaginatedResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

// First executes the query and returns only the first unreleased machine.
// Returns an error if no results are found.
//
// Example:
//
//	firstMachine, err := client.Machines.ListActive().
//		ByDifficulty("Insane").
//		First(ctx)
func (q *ActiveQuery) First(ctx context.Context) (MachinePaginatedResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return MachinePaginatedResponse{}, err
	}
	if len(resp.Data) == 0 {
		return MachinePaginatedResponse{}, errors.New("no results found")
	}
	return MachinePaginatedResponse{
		Data:         resp.Data[:1],
		ResponseMeta: resp.ResponseMeta,
	}, nil
}
