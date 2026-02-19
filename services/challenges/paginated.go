package challenges

import (
	"context"
	"errors"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
)

const (
	SortByChallenge   = v4Client.GetChallengesParamsSortBy("Challenge")
	SortByCategory    = v4Client.GetChallengesParamsSortBy("Catagory")
	SortByRating      = v4Client.GetChallengesParamsSortBy("Rating")
	SortByUsersSolves = v4Client.GetChallengesParamsSortBy("Solves")
	SortByReleaseDate = v4Client.GetChallengesParamsSortBy("ReleaseDate")
)

const (
	StateActive     = "active"
	StateRetired    = "retired"
	StateUnreleased = "unreleased"
)

type ChallengeListResponse struct {
	Data         []ChallengeList
	ResponseMeta common.ResponseMeta
}

// ByState filters challenges by state.
// Valid values are "active", "retired", and "unreleased".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByState("active").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Active challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByState(val string) *ChallengeQuery {
	return q.ByStateList(val)
}

// ByStateList filters challenges by multiple states.
// Valid values are "active", "retired", and "unreleased".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByStateList("active", "retired").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenges found: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByStateList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.state = lowercased
	return qc
}

// ByDifficulty filters challenges by difficulty level.
// Valid values are "VeryEasy", "Easy", "Medium", "Hard", and "Insane".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByDifficulty("Hard").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Hard challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByDifficulty(val string) *ChallengeQuery {
	return q.ByDifficultyList(val)
}

// ByDifficultyList filters challenges by multiple difficulty levels.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByDifficultyList("Hard", "Insane").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Hard/Insane challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByDifficultyList(val ...string) *ChallengeQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.difficulty = lowercased
	return qc
}

// ByCategory filters challenges by category ID.
// Category IDs are typically obtained from category listings.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByCategory(1).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Category matches: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByCategory(val ...int) *ChallengeQuery {
	return q.ByCategoryList(val...)
}

// ByCategoryList filters challenges by multiple category IDs.
// Category IDs are typically obtained from category listings.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByCategoryList(1, 2, 3).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Category matches: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByCategoryList(val ...int) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.category = val
	return qc
}

// SortedBy sets the field to sort results by.
// Valid values include "Challenge", "Category", "Rating", "Solves", and "ReleaseDate".
// Returns a new ChallengeQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	challenges, err := query.SortedBy("Rating").Descending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) SortedBy(field string) *ChallengeQuery {
	qc := ptr.Clone(q)
	sortBy := v4Client.GetChallengesParamsSortBy(field)
	qc.sortBy = sortBy
	return qc
}

func (q *ChallengeQuery) sort(val v4Client.GetChallengesParamsSortBy, order v4Client.GetChallengesParamsSortType) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.sortBy = val
	qc.sortType = order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.SortedBy("Rating").Ascending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) Ascending() *ChallengeQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetChallengesParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.SortedBy("Rating").Descending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) Descending() *ChallengeQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetChallengesParamsSortType("desc"))
}

// Page sets the specific page number for pagination.
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.Page(3).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Page 3 challenges: %d\n", len(challenges.Data))
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
//	challenges, err := query.PerPage(50).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenges in page: %d\n", len(challenges.Data))
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
//	challenges, err := query.Next().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Next page challenges: %d\n", len(challenges.Data))
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
//	challenges, err := query.Previous().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Previous page challenges: %d\n", len(challenges.Data))
func (q *ChallengeQuery) Previous() *ChallengeQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// ByKeyword filters Challenge name
// Returns a new ChallengeQuery that can be further chained.
//
// Example:
//
//	challenges, err := query.ByKeyword("spo").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Keyword matches: %d\n", len(challenges.Data))
func (q *ChallengeQuery) ByKeyword(keyword string) *ChallengeQuery {
	qc := ptr.Clone(q)
	qc.keyword = v4Client.Keyword(keyword)
	return qc
}

func (q *ChallengeQuery) fetchResults(ctx context.Context) (ChallengeListResponse, error) {
	params := &v4Client.GetChallengesParams{
		Page:    &q.page,
		PerPage: &q.perPage,
		Todo:    &q.todo,
	}

	if q.difficulty != nil {
		params.Difficulty = &q.difficulty
	}

	if q.category != nil {
		params.Category = &q.category
	}

	if q.state != nil {
		params.State = &q.state
	}

	if len(q.sortBy) > 0 {
		params.SortBy = &q.sortBy
	}

	if len(q.sortType) > 0 {
		params.SortType = &q.sortType
	}

	if len(q.status) > 0 {
		params.Status = &q.status
	}

	if q.keyword != "" {
		params.Keyword = &q.keyword
	}

	resp, err := q.client.V4().GetChallenges(q.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return ChallengeListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengesResponse)
	if err != nil {
		return ChallengeListResponse{ResponseMeta: meta}, err
	}

	return ChallengeListResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
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
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenges found: %d\n", len(challenges.Data))
func (q *ChallengeQuery) Results(ctx context.Context) (ChallengeListResponse, error) {
	return q.fetchResults(ctx)
}

type ChallengeList = v4Client.ChallengeList

// AllResults executes the query and returns all pages of challenges.
// This method automatically paginates through all available results.
// Use with caution for large datasets as it may consume significant memory.
//
// Example:
//
//	allChallenges, err := client.Challenges.List().
//		ByDifficulty("Hard").
//		AllResults(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Total challenges found: %d\n", len(allChallenges.Data))
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
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("First challenge: %s\n", firstChallenge.Data[0].Name)
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
