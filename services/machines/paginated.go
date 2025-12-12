package machines

import (
	"context"
	"errors"
	"strings"

	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"
)

type MachineQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        v5Client.GetMachinesParamsSortBy
	sortType      v5Client.GetMachinesParamsSortType
	difficulty    v5Client.MachineDifficulty
	os            v5Client.Os
	keyword       v5Client.Keyword
	state         v5Client.State
	free          *v5Client.GetMachinesParamsFree
	todo          *v5Client.GetMachinesParamsTodo
}

// List creates a new query for machines.
// This returns an MachineQuery that can be chained with filtering and pagination methods.
// Machines are the entities that can be queried with various filters.
//
// Example:
//
//	query := client.Machines.List()
//	machines, err := query.ByDifficulty("Hard").ByOS("Linux").Results(ctx)
func (s *Service) List() *MachineQuery {
	return &MachineQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

// Next moves to the next page in the pagination sequence.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	nextPage := query.Next().Results(ctx)
func (q *MachineQuery) Next() *MachineQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	prevPage := query.Previous().Results(ctx)
func (q *MachineQuery) Previous() *MachineQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// Page sets the specific page number for pagination.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.Page(3).Results(ctx)
func (q *MachineQuery) Page(n int) *MachineQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.PerPage(50).Results(ctx)
func (q *MachineQuery) PerPage(n int) *MachineQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// ByCompleted filters machines by completion status.
// Valid values are "Completed" and "InComplete".
// Returns a new MachineQuery that can be further chained.
// By default, it does not filter by completion status and
// includes both completed and incomplete machines.
//
// Example:
//
//	completed := query.ByCompleted("Completed").Results(ctx)
//	incomplete := query.ByCompleted("InComplete").Results(ctx)
func (q *MachineQuery) ByCompleted(val string) *MachineQuery {
	qc := ptr.Clone(q)
	qc.showCompleted = strings.ToLower(val)
	return qc
}

// ByOS filters machines by operating system.
// Valid values include "Linux" and "Windows".
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	linuxMachines := query.ByOS("Linux").Results(ctx)
//	linuxAndWindowsMachines := query.ByOS("Linux").ByOS("Windows").Results(ctx)
func (q *MachineQuery) ByOS(val string) *MachineQuery {
	return q.ByOSList(val)
}

// ByOSList filters machines by multiple operating systems.
// Valid values include "Linux" and "Windows".
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.ByOSList("Linux", "Windows").Results(ctx)
func (q *MachineQuery) ByOSList(val ...string) *MachineQuery {
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
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.ByDifficultyList("Hard", "Insane").Results(ctx)
func (q *MachineQuery) ByDifficultyList(val ...string) *MachineQuery {
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
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	hardMachines := query.ByDifficulty("Hard").Results(ctx)
//	mediumAndInsaneMachines := query.ByDifficulty("Medium").ByDifficulty("Insane").Results(ctx)
func (q *MachineQuery) ByDifficulty(val string) *MachineQuery {
	return q.ByDifficultyList(val)
}

func (q *MachineQuery) ByStateList(val ...string) *MachineQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.state = append(append([]string{}, q.state...), lowercased...)
	return qc
}

func (q *MachineQuery) ByState(val string) *MachineQuery {
	return q.ByStateList(val)
}

// SortedBy sets the field to sort results by.
// Valid values include "release-date", "name", "user-owns", "system-owns", "rating", "user-difficulty".
// Returns a new MachineQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	machines := query.SortedBy("name").Ascending().Results(ctx)
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *MachineQuery) SortedBy(field string) *MachineQuery {
	qc := ptr.Clone(q)
	sortBy := v5Client.GetMachinesParamsSortBy(strings.ToLower(field))
	qc.sortBy = sortBy
	return qc
}

func (q *MachineQuery) sort(val v5Client.GetMachinesParamsSortBy, order v5Client.GetMachinesParamsSortType) *MachineQuery {
	qc := ptr.Clone(q)
	qc.sortBy = val
	qc.sortType = order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Ascending().Results(ctx)
func (q *MachineQuery) Ascending() *MachineQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v5Client.GetMachinesParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *MachineQuery) Descending() *MachineQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v5Client.GetMachinesParamsSortType("desc"))
}

// Keyword filters machines names.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	machines := query.Keyword("buffer").Results(ctx)
func (q *MachineQuery) Keyword(val string) *MachineQuery {
	qc := ptr.Clone(q)
	qc.keyword = val
	return qc
}

// ByFree filters machines that are free or paid.
// This is useful for distinguishing between free and VIP machines.
// Also this has no bearing on active/retired status as there are free retired machines.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	freeMachines := query.ByFree(true).Results(ctx)
//	paidMachines := query.ByFree(false).Results(ctx)
func (q *MachineQuery) ByFree(val bool) *MachineQuery {
	v := v5Client.GetMachinesParamsFree(0)
	if val {
		v = 1
	}
	qc := ptr.Clone(q)
	qc.free = &v
	return qc
}

// ByTodo filters machines that have been marked as todo/favorite.
// Returns a new MachineQuery that can be further chained.
//
// Example:
//
//	todoMachines := query.ByTodo(true).Results(ctx)
//	nonTodoMachines := query.ByTodo(false).Results(ctx)
func (q *MachineQuery) ByTodo(val bool) *MachineQuery {
	v := v5Client.GetMachinesParamsTodo(0)
	if val {
		v = 1
	}
	qc := ptr.Clone(q)
	qc.todo = &v
	return qc
}

func (q *MachineQuery) fetchResults(ctx context.Context) (MachinesResponse, error) {
	params := &v5Client.GetMachinesParams{
		PerPage: &q.perPage,
		Page:    &q.page,
	}

	if q.keyword != "" {
		params.Keyword = &q.keyword
	}

	if len(q.difficulty) > 0 {
		d := q.difficulty
		params.Difficulty = &d
	}

	if len(q.os) > 0 {
		o := q.os
		params.Os = &o
	}

	if q.free != nil {
		params.Free = q.free
	}

	if q.todo != nil {
		params.Todo = q.todo
	}

	if q.showCompleted != "" {
		sc := v5Client.GetMachinesParamsShowCompleted(q.showCompleted)
		params.ShowCompleted = &sc
	}

	if len(q.state) > 0 {
		s := q.state
		params.State = &s
	}

	if q.sortBy != "" {
		sb := q.sortBy
		params.SortBy = &sb
	}

	if q.sortType != "" {
		st := q.sortType
		params.SortType = &st
	}

	resp, err := q.client.V5().GetMachines(q.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return MachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetMachinesResponse)
	if err != nil {
		return MachinesResponse{ResponseMeta: meta}, err
	}

	return MachinesResponse{
		Data:         wrapMachinesData(parsed.JSON200.Data),
		ResponseMeta: meta,
	}, nil
}

func wrapMachinesData(items []v5Client.MachinesItem) MachinesDataItems {
	out := make(MachinesDataItems, len(items))
	for i, item := range items {
		out[i] = MachinesData{MachinesItem: item}
	}
	return out
}

// Results executes the query and returns the current page of unreleased machines.
// This method should be called last in the query chain to fetch the actual data.
//
// Example:
//
//	machines, err := client.Machines.ListUnreleased().
//		ByDifficulty("Hard").
//		ByOS("Linux").
//		Page(1).
//		Results(ctx)
func (q *MachineQuery) Results(ctx context.Context) (MachinesResponse, error) {
	return q.fetchResults(ctx)
}

// AllResults executes the query and returns all pages of unreleased machines.
// This method automatically paginates through all available results.
// Use with caution for large datasets as it may consume significant memory.
//
// Example:
//
//	allMachines, err := client.Machines.ListUnreleased().
//		ByDifficulty("Hard").
//		AllResults(ctx)
func (q *MachineQuery) AllResults(ctx context.Context) (MachinesResponse, error) {
	var all MachinesDataItems
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return MachinesResponse{}, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return MachinesResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

// First executes the query and returns the first result of machines.
// If no results are found, an error is returned.
//
// Example:
//
//	firstMachine, err := client.Machines.List().
//		ByDifficulty("Hard").
//		First(ctx)
func (q *MachineQuery) First(ctx context.Context) (MachinesResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return MachinesResponse{}, err
	}
	if len(resp.Data) == 0 {
		return MachinesResponse{}, errors.New("no results found")
	}
	return MachinesResponse{
		Data:         resp.Data[:1],
		ResponseMeta: resp.ResponseMeta,
	}, nil
}
