package challenges

import (
	"context"
	"errors"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/ptr"
)

const (
	SortByChallenge   = v4Client.GetChallengesParamsSortBy("Challenge")
	SortByCatagory    = v4Client.GetChallengesParamsSortBy("Catagory")
	SortByRating      = v4Client.GetChallengesParamsSortBy("Rating")
	SortByUsersSolves = v4Client.GetChallengesParamsSortBy("Solves")
	SortByReleaseDate = v4Client.GetChallengesParamsSortBy("ReleaseDate")
)

const (
	StateActive     = "active"
	StateRetired    = "retired"
	StateUnreleased = "unreleased"
)

// ByState filters challenges by state.
// Valid values are "active", "retired", and "unreleased".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	activeChallenges := query.ByState("active").Results(ctx)
//	retiredChallenges := query.ByState("retired").Results(ctx)
func (q *ChallengeQuery) ByState(val string) *ChallengeQuery {
	return q.ByStateList(val)
}

// ByStateList filters challenges by multiple states.
// Valid values are "active", "retired", and "unreleased".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.ByStateList("active", "retired").Results(ctx)
func (q *ChallengeQuery) ByStateList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.state = &lowercased
	return qc
}

// ByDifficulty filters challenges by difficulty level.
// Valid values are "VeryEasy", "Easy", "Medium", "Hard", and "Insane".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	hardChallenges := query.ByDifficulty("Hard").Results(ctx)
//	easyChallenges := query.ByDifficulty("Hard").ByDifficulty("Easy").Results(ctx)
func (q *ChallengeQuery) ByDifficulty(val string) *ChallengeQuery {
	return q.ByDifficultyList(val)
}

// ByDifficultyList filters challenges by multiple difficulty levels.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.ByDifficultyList("Hard", "Insane").Results(ctx)
func (q *ChallengeQuery) ByDifficultyList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.difficulty = &lowercased
	return qc
}

// ByCategory filters challenges by category ID.
// Category IDs are typically obtained from category listings.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	webChallenges := query.ByCategory(1).Results(ctx)
//	cryptoChallenges := query.ByCategory(2).Results(ctx)
func (q *ChallengeQuery) ByCategory(val ...int) *ChallengeQuery {
	return q.ByCategoryList(val...)
}

// ByCategoryList filters challenges by multiple category IDs.
// Category IDs are typically obtained from category listings.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.ByCategoryList(1, 2, 3).Results(ctx)
func (q *ChallengeQuery) ByCategoryList(val ...int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.category = &val
	return qc
}

// SortedBy sets the field to sort results by.
// Valid values include "Challenge", "Catagory", "Rating", "Solves", and "ReleaseDate".
// Returns a new ChallengeQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	challenges := query.SortedBy("Rating").Descending().Results(ctx)
//	challenges := query.SortedBy("ReleaseDate").Ascending().Results(ctx)
func (q *ChallengeQuery) SortedBy(field string) *ChallengeQuery {
	qc := ptr.Clone(q)
	sortBy := v4Client.GetChallengesParamsSortBy(field)
	qc.sortBy = &sortBy
	return qc
}

func (q *ChallengeQuery) sort(val v4Client.GetChallengesParamsSortBy, order v4Client.GetChallengesParamsSortType) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.sortBy = &val
	qc.sortType = &order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.SortedBy("Rating").Ascending().Results(ctx)
func (q *ChallengeQuery) Ascending() *ChallengeQuery {
	if q.sortBy == nil {
		return q
	}
	return q.sort(*q.sortBy, v4Client.GetChallengesParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.SortedBy("Rating").Descending().Results(ctx)
func (q *ChallengeQuery) Descending() *ChallengeQuery {
	if q.sortBy == nil {
		return q
	}
	return q.sort(*q.sortBy, v4Client.GetChallengesParamsSortType("desc"))
}

// Page sets the specific page number for pagination.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.Page(3).Results(ctx)
func (q *ChallengeQuery) Page(n int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges := query.PerPage(50).Results(ctx)
func (q *ChallengeQuery) PerPage(n int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// Next moves to the next page in the pagination sequence.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	nextPage := query.Next().Results(ctx)
func (q *ChallengeQuery) Next() *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	prevPage := query.Previous().Results(ctx)
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

// Results executes the query and returns the current page of challenges.
// This method should be called last in the query chain to fetch the actual data.
//
// Example:
//
//	challenges, err := client.Challenges.List().
//		ByDifficulty("Hard").
//		ByState("active").
//		SortedBy("Rating").Descending().
//		Results(ctx)
func (q *ChallengeQuery) Results(ctx context.Context) (ChallengeListResponse, error) {
	return q.fetchResults(ctx)
}

// AllResults executes the query and returns all pages of challenges.
// This method automatically paginates through all available results.
// Use with caution for large datasets as it may consume significant memory.
//
// Example:
//
//	allChallenges, err := client.Challenges.List().
//		ByDifficulty("Hard").
//		AllResults(ctx)
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

// First executes the query and returns only the first challenge.
// Returns an error if no results are found.
//
// Example:
//
//	firstChallenge, err := client.Challenges.List().
//		ByDifficulty("Insane").
//		SortedBy("Rating").Descending().
//		First(ctx)
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
