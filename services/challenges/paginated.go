package challenges

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

const (
	SortByChallenge   = v4Client.GetChallengesParamsSortBy("Challenge")
	SortByCatagory    = v4Client.GetChallengesParamsSortBy("Catagory")
	SortByRating      = v4Client.GetChallengesParamsSortBy("Rating")
	SortByUsersSolves = v4Client.GetChallengesParamsSortBy("UsersSolves")
	SortByReleaseDate = v4Client.GetChallengesParamsSortBy("ReleaseDate")
)

const (
	StateActive     = "active"
	StateRetired    = "retired"
	StateUnreleased = "unreleased"
)

func (q *ChallengeQuery) ByState(val string) *ChallengeQuery {
	return q.ByStateList(val)
}

func (q *ChallengeQuery) ByStateList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.state = &val
	return qc
}

func (q *ChallengeQuery) ByDifficulty(val string) *ChallengeQuery {
	return q.ByDifficultyList(val)
}

func (q *ChallengeQuery) ByDifficultyList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.difficulty = &val
	return qc
}

func (q *ChallengeQuery) ByCategory(val ...int) *ChallengeQuery {
	return q.ByCategoryList(val...)
}

func (q *ChallengeQuery) ByCategoryList(val ...int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.category = &val
	return qc
}

func (q *ChallengeQuery) SortedBy(val v4Client.GetChallengesParamsSortBy) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	return qc
}

func (q *ChallengeQuery) sort(val v4Client.GetChallengesParamsSortBy, order v4Client.GetChallengesParamsSortType) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	qc.sortType = &order
	return qc
}

func (q *ChallengeQuery) Ascending() *ChallengeQuery {
	if q.sortBy == nil {
		return q // or panic/log if you want to enforce setting sortBy first
	}
	return q.sort(*q.sortBy, v4Client.GetChallengesParamsSortType("asc"))
}

func (q *ChallengeQuery) Descending() *ChallengeQuery {
	if q.sortBy == nil {
		return q
	}
	return q.sort(*q.sortBy, v4Client.GetChallengesParamsSortType("desc"))
}

func (q *ChallengeQuery) Page(n int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

func (q *ChallengeQuery) PerPage(n int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

func (q *ChallengeQuery) Next() *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

func (q *ChallengeQuery) Previous() *ChallengeQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

func (q *ChallengeQuery) fetchResults(ctx context.Context) (ChallengeListResponse, error) {
	params := &v4Client.GetChallengesParams{
		Page:       &q.page,
		PerPage:    &q.perPage,
		Difficulty: q.difficulty,
		Category:   q.category,
		SortBy:     q.sortBy,
		Status:     q.status,
		SortType:   q.sortType,
		State:      q.state,
		Todo:       q.todo,
	}

	resp, err := q.client.V4().GetChallengesWithResponse(q.client.Limiter().Wrap(ctx), params)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ChallengeListResponse {
			return ChallengeListResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ChallengeListResponse{
		Data: convert.Slice(*resp.JSON200.Data, fromAPIChallengeList),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
			CFRay:      resp.HTTPResponse.Header.Get("CF-Ray"),
		},
	}, nil
}

func (q *ChallengeQuery) Results(ctx context.Context) (ChallengeListResponse, error) {
	return q.fetchResults(ctx)
}

func (q *ChallengeQuery) AllResults(ctx context.Context) (ChallengeListResponse, error) {
	var all []ChallengeList
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return ChallengeListResponse{}, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return ChallengeListResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

func (q *ChallengeQuery) First(ctx context.Context) (ChallengeListResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return ChallengeListResponse{}, err
	}
	if len(resp.Data) == 0 {
		return ChallengeListResponse{}, errors.New("no results found")
	}
	return ChallengeListResponse{
		Data:         resp.Data[:1],
		ResponseMeta: resp.ResponseMeta,
	}, nil
}
