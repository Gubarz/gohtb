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

func (s *Service) ListUnreleased() *UnreleasedQuery {
	return &UnreleasedQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

func (q *UnreleasedQuery) Next() *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

func (q *UnreleasedQuery) Previous() *UnreleasedQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

func (q *UnreleasedQuery) Page(n int) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

func (q *UnreleasedQuery) PerPage(n int) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// Completed, InComplete
func (q *UnreleasedQuery) ByCompleted(val string) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.showCompleted = val
	return qc
}

// Linux, Windows
func (q *UnreleasedQuery) ByOS(val string) *UnreleasedQuery {
	return q.ByOSList(val)
}

// Linux, Windows
func (q *UnreleasedQuery) ByOSList(val ...string) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.os = append(append([]string{}, q.os...), val...)
	return qc
}

// Easy, Medium, Hard, Insane
func (q *UnreleasedQuery) ByDifficultyList(val ...string) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.difficulty = append(append([]string{}, q.difficulty...), val...)
	return qc
}

// Easy, Medium, Hard, Insane
func (q *UnreleasedQuery) ByDifficulty(val string) *UnreleasedQuery {
	return q.ByDifficultyList(val)
}

func (q *UnreleasedQuery) SortedBy(val v4Client.GetMachinePaginatedParamsSortBy) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	return qc
}

func (q *UnreleasedQuery) sort(val v4Client.GetMachinePaginatedParamsSortBy, order v4Client.GetMachinePaginatedParamsSortType) *UnreleasedQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	qc.sortType = &order
	return qc
}

func (q *UnreleasedQuery) Ascending() *UnreleasedQuery {
	if q.sortBy == nil {
		return q // or panic/log if you want to enforce setting sortBy first
	}
	return q.sort(*q.sortBy, v4Client.GetMachinePaginatedParamsSortType("asc"))
}

func (q *UnreleasedQuery) Descending() *UnreleasedQuery {
	if q.sortBy == nil {
		return q
	}
	return q.sort(*q.sortBy, v4Client.GetMachinePaginatedParamsSortType("desc"))
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
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) MachineUnreleasedResponse {
			return MachineUnreleasedResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MachineUnreleasedResponse{
		Data: convert.Slice(*resp.JSON200.Data, fromAPIMachineUnreleasedData),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (q *UnreleasedQuery) Results(ctx context.Context) (MachineUnreleasedResponse, error) {
	return q.fetchResults(ctx)
}

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
