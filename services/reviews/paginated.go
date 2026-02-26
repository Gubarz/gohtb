package reviews

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"
)

// ReviewQuery builds and executes paginated review queries.
type ReviewQuery struct {
	client    service.Client
	product   Product
	productId int
	page      int
	perPage   int
}

// Next moves to the next page in the pagination sequence.
// Returns a new ReviewQuery that can be further chained.
//
// Example:
//
//	reviewsPage, err := query.Next().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Next page review entries: %d\n", len(reviewsPage.Data.Data))
func (q *ReviewQuery) Next() *ReviewQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

// Previous moves to the previous page in the pagination sequence.
// If already on the first page, it remains on page 1.
// Returns a new ReviewQuery that can be further chained.
//
// Example:
//
//	reviewsPage, err := query.Previous().Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Previous page review entries: %d\n", len(reviewsPage.Data.Data))
func (q *ReviewQuery) Previous() *ReviewQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

// Page sets the specific page number for pagination.
// Returns a new ReviewQuery that can be further chained.
//
// Example:
//
//	reviewsPage, err := query.Page(3).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Page 3 review entries: %d\n", len(reviewsPage.Data.Data))
func (q *ReviewQuery) Page(n int) *ReviewQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

// PerPage sets the number of results per page.
// Returns a new ReviewQuery that can be further chained.
//
// Example:
//
//	reviewsPage, err := query.PerPage(25).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review entries in page: %d\n", len(reviewsPage.Data.Data))
func (q *ReviewQuery) PerPage(n int) *ReviewQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

func (q *ReviewQuery) fetchResults(ctx context.Context) (ReviewPaginatedResponse, error) {
	paramsEditor := func(_ context.Context, req *http.Request) error {
		query := req.URL.Query()
		if q.page > 0 {
			query.Set("page", strconv.Itoa(q.page))
		}
		if q.perPage > 0 {
			query.Set("per_page", strconv.Itoa(q.perPage))
		}
		req.URL.RawQuery = query.Encode()
		return nil
	}

	resp, err := q.client.V4().GetReviewPaginated(
		q.client.Limiter().Wrap(ctx),
		v4Client.GetReviewPaginatedParamsProduct(q.product),
		q.productId,
		paramsEditor,
	)
	if err != nil {
		return ReviewPaginatedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetReviewPaginatedResponse)
	if err != nil {
		return ReviewPaginatedResponse{ResponseMeta: meta}, err
	}

	return ReviewPaginatedResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// Results executes the query and returns the current page of reviews.
// This method should be called last in the query chain to fetch the actual data.
//
// Example:
//
//	reviewsPage, err := client.Reviews.Machine(12345).List().
//		PerPage(25).
//		Page(1).
//		Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review entries found: %d\n", len(reviewsPage.Data.Data))
func (q *ReviewQuery) Results(ctx context.Context) (ReviewPaginatedResponse, error) {
	return q.fetchResults(ctx)
}

// AllResults executes the query and returns all pages of reviews.
// This method automatically paginates through all available results.
//
// Example:
//
//	allReviews, err := client.Reviews.Machine(12345).List().AllResults(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Total review entries found: %d\n", len(allReviews.Data.Data))
func (q *ReviewQuery) AllResults(ctx context.Context) (ReviewPaginatedResponse, error) {
	var all v4Client.ReviewMessage
	page := q.page
	var merged ReviewPaginatedResponse

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return ReviewPaginatedResponse{}, err
		}

		if page == q.page {
			merged = resp
		}

		all = append(all, resp.Data.Data...)
		merged.ResponseMeta = resp.ResponseMeta

		if resp.Data.Meta.LastPage > 0 && page >= resp.Data.Meta.LastPage {
			break
		}

		if len(resp.Data.Data) == 0 {
			break
		}

		if q.perPage > 0 && len(resp.Data.Data) < q.perPage {
			break
		}

		page++
	}

	merged.Data.Data = all
	merged.Data.Count = len(all)

	return merged, nil
}

// First executes the query and returns only the first review entry.
// Returns an error if no results are found.
//
// Example:
//
//	firstReview, err := client.Reviews.Machine(12345).List().First(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("First review ID: %d\n", firstReview.Data.Data[0].Id)
func (q *ReviewQuery) First(ctx context.Context) (ReviewPaginatedResponse, error) {
	resp, err := q.fetchResults(ctx)
	if err != nil {
		return ReviewPaginatedResponse{}, err
	}

	if len(resp.Data.Data) == 0 {
		return ReviewPaginatedResponse{}, errors.New("no results found")
	}

	resp.Data.Data = resp.Data.Data[:1]
	resp.Data.Count = 1
	return resp, nil
}
