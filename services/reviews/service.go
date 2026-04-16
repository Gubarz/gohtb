package reviews

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to product review endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new reviews service bound to a shared client.
//
// Example:
//
//	reviewsService := reviews.NewService(client)
//	_ = reviewsService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

// Product identifies a reviewable product type.
type Product = string

const (
	ProductChallenge Product = "challenge"
	ProductMachine   Product = "machine"
	ProductSherlock  Product = "sherlock"
)

// Handle scopes review operations to a specific product and product ID.
type Handle struct {
	client    service.Client
	product   Product
	productId int
}

type ReviewPaginatedData = v4Client.ReviewProductPaginatedResponse

// ReviewPaginatedResponse contains paginated review entries for a product.
type ReviewPaginatedResponse struct {
	Data         ReviewPaginatedData
	ResponseMeta common.ResponseMeta
}

// Target returns a review handle for the specified product and product ID.
//
// Example:
//
//	reviewTarget := client.Reviews.Target(reviews.ProductMachine, 12345)
//	_ = reviewTarget
func (s *Service) Target(product Product, productID int) *Handle {
	return &Handle{
		client:    s.base.Client,
		product:   product,
		productId: productID,
	}
}

// Challenge returns a review handle for a challenge.
//
// Example:
//
//	challengeReviews := client.Reviews.Challenge(12345)
//	_ = challengeReviews
func (s *Service) Challenge(challengeId int) *Handle {
	return s.Target(ProductChallenge, challengeId)
}

// Machine returns a review handle for a machine.
//
// Example:
//
//	machineReviews := client.Reviews.Machine(12345)
//	_ = machineReviews
func (s *Service) Machine(machineId int) *Handle {
	return s.Target(ProductMachine, machineId)
}

// Sherlock returns a review handle for a sherlock.
//
// Example:
//
//	sherlockReviews := client.Reviews.Sherlock(12345)
//	_ = sherlockReviews
func (s *Service) Sherlock(sherlockId int) *Handle {
	return s.Target(ProductSherlock, sherlockId)
}

// List creates a paginated review query for the handle target.
//
// Example:
//
//	query := client.Reviews.Machine(12345).List()
//	_ = query
func (h *Handle) List() *ReviewQuery {
	return &ReviewQuery{
		client:    h.client,
		product:   h.product,
		productId: h.productId,
		page:      1,
		perPage:   100,
	}
}

// Paginated retrieves the first page of reviews for the handle target.
//
// Example:
//
//	reviewsPage, err := client.Reviews.Challenge(1).Paginated(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review entries on page: %d\n", len(reviewsPage.Data.Data))
func (h *Handle) Paginated(ctx context.Context) (ReviewPaginatedResponse, error) {
	return h.List().Results(ctx)
}
