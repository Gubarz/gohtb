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

func (s *Service) ListRetired() *RetiredQuery {
	return &RetiredQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

func (q *RetiredQuery) Next() *RetiredQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

func (q *RetiredQuery) Previous() *RetiredQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

func (q *RetiredQuery) Page(n int) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

func (q *RetiredQuery) PerPage(n int) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// Completed, InComplete
func (q *RetiredQuery) ByCompleted(val string) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.showCompleted = val
	return qc
}

// Linux, Windows
func (q *RetiredQuery) ByOS(val string) *RetiredQuery {
	return q.ByOSList(val)
}

// Linux, Windows
func (q *RetiredQuery) ByOSList(val ...string) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.os = append(append([]string{}, q.os...), val...)
	return qc
}

// Easy, Medium, Hard, Insane
func (q *RetiredQuery) ByDifficultyList(val ...string) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.difficulty = append(append([]string{}, q.difficulty...), val...)
	return qc
}

// Easy, Medium, Hard, Insane
func (q *RetiredQuery) ByDifficulty(val string) *RetiredQuery {
	return q.ByDifficultyList(val)
}

func (q *RetiredQuery) SortedBy(val v4Client.GetMachinePaginatedParamsSortBy) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	return qc
}

func (q *RetiredQuery) sort(val v4Client.GetMachinePaginatedParamsSortBy, order v4Client.GetMachinePaginatedParamsSortType) *RetiredQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	qc.sortType = &order
	return qc
}

func (q *RetiredQuery) Ascending() *RetiredQuery {
	if q.sortBy == nil {
		return q // or panic/log if you want to enforce setting sortBy first
	}
	return q.sort(*q.sortBy, v4Client.GetMachinePaginatedParamsSortType("asc"))
}

func (q *RetiredQuery) Descending() *RetiredQuery {
	if q.sortBy == nil {
		return q
	}
	return q.sort(*q.sortBy, v4Client.GetMachinePaginatedParamsSortType("desc"))
}

func (q *RetiredQuery) fetchResults(ctx context.Context) (MachinePaginatedResponse, error) {
	params := &v4Client.GetMachineListRetiredPaginatedParams{
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

	if len(q.tags) > 0 {
		d := q.tags
		params.Tags = &d
	}

	if q.showCompleted != "" {
		sc := v4Client.GetMachineListRetiredPaginatedParamsShowCompleted(q.showCompleted)
		params.ShowCompleted = &sc
	}

	resp, err := q.client.V4().GetMachineListRetiredPaginatedWithResponse(q.client.Limiter().Wrap(ctx), params)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) MachinePaginatedResponse {
			return MachinePaginatedResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MachinePaginatedResponse{
		Data: convert.Slice(*resp.JSON200.Data, fromAPIMachineData),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (q *RetiredQuery) Results(ctx context.Context) (MachinePaginatedResponse, error) {
	return q.fetchResults(ctx)
}

func (q *RetiredQuery) AllResults(ctx context.Context) (MachinePaginatedResponse, error) {
	var all []MachineData
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return MachinePaginatedResponse{}, err
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

func (q *RetiredQuery) First(ctx context.Context) (MachinePaginatedResponse, error) {
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
