package sherlocks

import (
	"context"
	"errors"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
)

const (
	SortBySherlock    = v4Client.GetSherlocksParamsSortBy("Sherlock")
	SortByCategory    = v4Client.GetSherlocksParamsSortBy("Catagory")
	SortByRating      = v4Client.GetSherlocksParamsSortBy("Rating")
	SortByUsersSolves = v4Client.GetSherlocksParamsSortBy("Solves")
	SortByReleaseDate = v4Client.GetSherlocksParamsSortBy("ReleaseDate")
)

const (
	StateActive     = "active"
	StateRetired    = "retired"
	StateUnreleased = "unreleased"
)

type SherlockItem = v4Client.SherlockItem

type SherlockListResponse struct {
	Data         []SherlockItem
	ResponseMeta common.ResponseMeta
}

// ByState filters Sherlocks by state.
// Valid values are "active", "retired", and "unreleased".
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByState("active").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Active sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByState(val string) *SherlockQuery {
	return q.ByStateList(val)
}

// ByStateList filters Sherlocks by multiple states.
// Valid values are "active", "retired", and "unreleased".
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByStateList("active", "retired").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlocks found: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByStateList(val ...string) *SherlockQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.state = lowercased
	return qc
}

// ByDifficulty filters Sherlocks by difficulty level.
// Valid values are "VeryEasy", "Easy", "Medium", "Hard", and "Insane".
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByDifficulty("Hard").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Hard sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByDifficulty(val string) *SherlockQuery {
	return q.ByDifficultyList(val)
}

// ByDifficultyList filters Sherlocks by multiple difficulty levels.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByDifficultyList("Hard", "Insane").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Hard/Insane sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByDifficultyList(val ...string) *SherlockQuery {
	qc := ptr.Clone(q)
	lowercased := make([]string, len(val))
	for i, v := range val {
		lowercased[i] = strings.ToLower(v)
	}
	qc.difficulty = lowercased
	return qc
}

// ByCategory filters Sherlocks by category ID.
// Category IDs are typically obtained from category listings.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByCategory(1).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Category matches: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByCategory(val ...int) *SherlockQuery {
	return q.ByCategoryList(val...)
}

// ByCategoryList filters Sherlocks by multiple category IDs.
// Category IDs are typically obtained from category listings.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByCategoryList(1, 2, 3).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Category matches: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByCategoryList(val ...int) *SherlockQuery {
	qc := ptr.Clone(q)
	qc.category = val
	return qc
}

// SortedBy sets the field to sort results by.
// Valid values include "Sherlock", "Category", "Rating", "Solves", and "ReleaseDate".
// Returns a new SherlockQuery that can be further chained with Ascending() or Descending().
//
// Example:
//
//	sherlocks, err := query.SortedBy("Rating").Descending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) SortedBy(field string) *SherlockQuery {
	qc := ptr.Clone(q)
	sortBy := v4Client.GetSherlocksParamsSortBy(field)
	qc.sortBy = sortBy
	return qc
}

func (q *SherlockQuery) sort(val v4Client.GetSherlocksParamsSortBy, order v4Client.GetSherlocksParamsSortType) *SherlockQuery {
	qc := ptr.Clone(q)
	qc.sortBy = val
	qc.sortType = order
	return qc
}

// Ascending sets the sort order to ascending.
// Must be called after SortedBy(). Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.SortedBy("Rating").Ascending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Ascending() *SherlockQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetSherlocksParamsSortType("asc"))
}

// Descending sets the sort order to descending.
// Must be called after SortedBy(). Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.SortedBy("Rating").Descending().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sorted sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Descending() *SherlockQuery {
	if q.sortBy == "" {
		return q
	}
	return q.sort(q.sortBy, v4Client.GetSherlocksParamsSortType("desc"))
}

// Page sets the specific page number for pagination.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.Page(3).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Page 3 sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Page(n int) *SherlockQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.PerPage(50).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlocks in page: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) PerPage(n int) *SherlockQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

// Next moves to the next page in the pagination sequence.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.Next().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Next page sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Next() *SherlockQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.Previous().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Previous page sherlocks: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Previous() *SherlockQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// ByKeyword filters Sherlocks name
// Returns a new SherlockQuery that can be further chained.
//
// Example:
//
//	sherlocks, err := query.ByKeyword("bru").Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Keyword matches: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) ByKeyword(keyword string) *SherlockQuery {
	qc := ptr.Clone(q)
	qc.keyword = v4Client.Keyword(keyword)
	return qc
}

func (q *SherlockQuery) fetchResults(ctx context.Context) (SherlockListResponse, error) {
	params := &v4Client.GetSherlocksParams{
		Page:    &q.page,
		PerPage: &q.perPage,
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

	resp, err := q.client.V4().GetSherlocks(q.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return SherlockListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlocksResponse)
	if err != nil {
		return SherlockListResponse{ResponseMeta: meta}, err
	}

	return SherlockListResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// Results executes the query and returns the current page of Sherlocks.
// This method should be called last in the query chain to fetch the actual data.
//
// Example:
//
//	sherlocks, err := client.Sherlocks.List().
//		ByDifficulty("Hard").
//		ByState("active").
//		SortedBy("Rating").Descending().
//		Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlocks found: %d\n", len(sherlocks.Data))
func (q *SherlockQuery) Results(ctx context.Context) (SherlockListResponse, error) {
	return q.fetchResults(ctx)
}

// AllResults executes the query and returns all pages of Sherlocks.
// This method automatically paginates through all available results.
// Use with caution for large datasets as it may consume significant memory.
//
// Example:
//
//	allSherlocks, err := client.Sherlocks.List().
//		ByDifficulty("Hard").
//		AllResults(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Total sherlocks found: %d\n", len(allSherlocks.Data))
func (q *SherlockQuery) AllResults(ctx context.Context) (SherlockListResponse, error) {
	var all []SherlockItem
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return SherlockListResponse{}, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return SherlockListResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

// First executes the query and returns only the first Sherlock.
// Returns an error if no results are found.
//
// Example:
//
//	firstSherlock, err := client.Sherlocks.List().
//		ByDifficulty("Insane").
//		SortedBy("Rating").Descending().
//		First(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("First sherlock: %s\n", firstSherlock.Data[0].Name)
func (q *SherlockQuery) First(ctx context.Context) (SherlockListResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return SherlockListResponse{}, err
	}
	if len(resp.Data) == 0 {
		return SherlockListResponse{}, errors.New("no results found")
	}
	return SherlockListResponse{
		Data:         resp.Data[:1],
		ResponseMeta: resp.ResponseMeta,
	}, nil
}
