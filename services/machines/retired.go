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

type RetiredQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        v5Client.GetMachinesParamsSortBy
	sortType      v5Client.GetMachinesParamsSortType
	difficulty    v5Client.MachineDifficulty
	os            v5Client.Os
	keyword       v5Client.Keyword
	// tags          []v5Client.MachineLabel v5 has no tags
}

// ListRetired creates a new query for retired machines.
// This returns an RetiredQuery that can be chained with filtering and pagination methods.
// Retired machines are machines that are no longer publicly available.
//
// Note: Note: Updated by Ceald 🐱
//
// Example:
//
//	query := client.Machines.ListRetired()
//	machines, err := query.ByDifficulty("Hard").ByOS("Linux").Results(ctx)
func (s *Service) ListRetired() *RetiredQuery {
	return &RetiredQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

// Next moves to the next page in the pagination sequence.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	nextPage := query.Next().Results(ctx)
func (q *RetiredQuery) Next() *RetiredQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	prevPage := query.Previous().Results(ctx)
func (q *RetiredQuery) Previous() *RetiredQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// Page sets the specific page number for pagination.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.Page(3).Results(ctx)
func (q *RetiredQuery) Page(n int) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.PerPage(50).Results(ctx)
func (q *RetiredQuery) PerPage(n int) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// ByCompleted filters machines by completion status.
// Valid values are "Completed" and "InComplete".
// Returns a new RetiredQuery that can be further chained.
// By default, it does not filter by completion status and
// includes both completed and incomplete machines.
//
// Example:
//
//	completed := query.ByCompleted("Completed").Results(ctx)
//	incomplete := query.ByCompleted("InComplete").Results(ctx)
func (q *RetiredQuery) ByCompleted(val string) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.showCompleted = strings.ToLower(val)
	return qc
}

// ByOS filters machines by operating system.
// Valid values include "Linux" and "Windows".
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	linuxMachines := query.ByOS("Linux").Results(ctx)
//	linuxAndWindowsMachines := query.ByOS("Linux").ByOS("Windows").Results(ctx)
func (q *RetiredQuery) ByOS(val string) *RetiredQuery {
	return q.ByOSList(val)
}

// ByOSList filters machines by multiple operating systems.
// Valid values include "Linux" and "Windows".
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.ByOSList("Linux", "Windows").Results(ctx)
func (q *RetiredQuery) ByOSList(val ...string) *RetiredQuery {
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
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.ByDifficultyList("Hard", "Insane").Results(ctx)
func (q *RetiredQuery) ByDifficultyList(val ...string) *RetiredQuery {
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
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	hardMachines := query.ByDifficulty("Hard").Results(ctx)
//	mediumAndInsaneMachines := query.ByDifficulty("Medium").ByDifficulty("Insane").Results(ctx)
func (q *RetiredQuery) ByDifficulty(val string) *RetiredQuery {
	return q.ByDifficultyList(val)
}

// SortedBy sets the field to sort results by.
// Valid values include "release-date", "name", "user-owns", "system-owns", "rating", "user-difficulty".
// Returns a new RetiredQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	machines := query.SortedBy("name").Ascending().Results(ctx)
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *RetiredQuery) SortedBy(field string) *RetiredQuery {
	qc := ptr.Clone(q)
	sortBy := v5Client.GetMachinesParamsSortBy(strings.ToLower(field))
	qc.sortBy = sortBy
	return qc
}

func (q *RetiredQuery) sort(val v5Client.GetMachinesParamsSortBy, order v5Client.GetMachinesParamsSortType) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.sortBy = val
	qc.sortType = order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Ascending().Results(ctx)
func (q *RetiredQuery) Ascending() *RetiredQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v5Client.GetMachinesParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.SortedBy("user-difficulty").Descending().Results(ctx)
func (q *RetiredQuery) Descending() *RetiredQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v5Client.GetMachinesParamsSortType("desc"))
}


// ** Tags are obsolete in v5

// Keyword filters machines names.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.Keyword("buffer").Results(ctx)
// func (q *RetiredQuery) Keyword(val string) *RetiredQuery {
// 	qc := ptr.Clone(q)
// 	qc.keyword = val
// 	return qc
// }

// // ByTag filters machines by a single tag ID.
// // Returns a new RetiredQuery that can be further chained.
// //
// // Example:
// //
// //	machines := query.ByTag(123).Results(ctx)
// func (q *RetiredQuery) ByTag(val int) *RetiredQuery {
// 	return q.ByTagList(val)
// }

// ByTagList filters machines by multiple tag IDs.
// Returns a new RetiredQuery that can be further chained.
//
// Example:
//
//	machines := query.ByTagList(123, 456).Results(ctx)
// func (q *RetiredQuery) ByTagList(val ...string) *RetiredQuery {
// 	qc := ptr.Clone(q)
// 	qc.tags = append(append([]string{}, q.tags...), val...)
// 	return qc
// }

// **



func (q *RetiredQuery) fetchResults(ctx context.Context) (MachinesResponse, error) {
	var state v5Client.State
	state = append(state, "retired")
	params := &v5Client.GetMachinesParams{
		PerPage: &q.perPage,
		Page:    &q.page,
		Keyword: &q.keyword,
		State: &state,

	}

	if len(q.difficulty) > 0 {
		d := q.difficulty
		params.Difficulty = &d
	}

	if len(q.os) > 0 {
		o := q.os
		params.Os = &o
	}

	// if len(q.tags) > 0 {
	// 	d := q.tags
	// 	params.Tag = &d
	// }

	if q.showCompleted != "" {
		sc := v5Client.GetMachinesParamsShowCompleted(q.showCompleted)
		params.ShowCompleted = &sc
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
func (q *RetiredQuery) Results(ctx context.Context) (MachinesResponse, error) {
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
func (q *RetiredQuery) AllResults(ctx context.Context) (MachinesResponse, error) {
	var all []MachinesData
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

func (q *RetiredQuery) First(ctx context.Context) (MachinesResponse, error) {
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
