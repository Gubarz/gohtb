package machines

import (
	"context"
	"errors"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v4Client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/ptr"
)

// ListUnreleased creates a new query for unreleased machines.
// This returns an UnreleasedQuery that can be chained with filtering and pagination methods.
// Unreleased machines are machines that are not yet publicly available.
//
// Example:
//
//	query := client.Machines.ListUnreleased()
//	machines, err := query.ByDifficulty("Hard").ByOS("Linux").Results(ctx)
func (s *Service) ListUnreleased() *UnreleasedQuery {
	return &UnreleasedQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

// Next moves to the next page in the pagination sequence.
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	nextPage := query.Next().Results(ctx)
func (q *UnreleasedQuery) Next() *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	prevPage := query.Previous().Results(ctx)
func (q *UnreleasedQuery) Previous() *UnreleasedQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// Page sets the specific page number for pagination.
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	machines := query.Page(3).Results(ctx)
func (q *UnreleasedQuery) Page(n int) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	machines := query.PerPage(50).Results(ctx)
func (q *UnreleasedQuery) PerPage(n int) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// ByOS filters machines by operating system.
// Valid values include "Linux" and "Windows".
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	linuxMachines := query.ByOS("Linux").Results(ctx)
//	linuxAndWindowsMachines := query.ByOS("Linux").ByOS("Windows").Results(ctx)
func (q *UnreleasedQuery) ByOS(val string) *UnreleasedQuery {
	return q.ByOSList(val)
}

// ByOSList filters machines by multiple operating systems.
// Valid values include "Linux" and "Windows".
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	machines := query.ByOSList("Linux", "Windows").Results(ctx)
func (q *UnreleasedQuery) ByOSList(val ...string) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.os = append(append([]string{}, q.os...), val...)
	return qc
}

// ByDifficulty filters machines by difficulty level.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	hardMachines := query.ByDifficulty("Hard").Results(ctx)
//	mediumAndInsaneMachines := query.ByDifficulty("Medium").ByDifficulty("Insane").Results(ctx)
func (q *UnreleasedQuery) ByDifficultyList(val ...string) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.difficulty = append(append([]string{}, q.difficulty...), val...)
	return qc
}

// ByDifficultyList filters machines by multiple difficulty levels.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new UnreleasedQuery that can be further chained.
//
// Example:
//
//	machines := query.ByDifficultyList("Hard", "Insane").Results(ctx)
func (q *UnreleasedQuery) ByDifficulty(val string) *UnreleasedQuery {
	return q.ByDifficultyList(val)
}

func (q *UnreleasedQuery) fetchResults(ctx context.Context) (MachineUnreleasedResponse, error) {
	params := &v4Client.GetMachineUnreleasedParams{
		PerPage: &q.perPage,
		Page:    &q.page,
		Keyword: q.keyword,
	}

	if len(q.difficulty) > 0 {
		d := q.difficulty
		params.Difficulty = &d
	}

	if len(q.os) > 0 {
		o := q.os
		params.Os = &o
	}

	resp, err := q.client.V4().GetMachineUnreleasedWithResponse(q.client.Limiter().Wrap(ctx), params)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) MachineUnreleasedResponse {
			return MachineUnreleasedResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MachineUnreleasedResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPIMachineUnreleasedData),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
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
func (q *UnreleasedQuery) Results(ctx context.Context) (MachineUnreleasedResponse, error) {
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
func (q *UnreleasedQuery) AllResults(ctx context.Context) (MachineUnreleasedResponse, error) {
	var all []MachineUnreleasedData
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return MachineUnreleasedResponse{}, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return MachineUnreleasedResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

// First executes the query and returns only the first unreleased machine.
// Returns an error if no results are found.
//
// Example:
//
//	firstMachine, err := client.Machines.ListUnreleased().
//		ByDifficulty("Insane").
//		First(ctx)
func (q *UnreleasedQuery) First(ctx context.Context) (MachineUnreleasedResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return MachineUnreleasedResponse{}, err
	}
	if len(resp.Data) == 0 {
		return MachineUnreleasedResponse{}, errors.New("no results found")
	}
	return MachineUnreleasedResponse{
		Data:         resp.Data[:1],
		ResponseMeta: resp.ResponseMeta,
	}, nil
}
