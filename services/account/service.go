package account

import (
	"context"
	"errors"

	experienceClient "github.com/gubarz/gohtb/httpclient/experience"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

// NewService creates a new ids service bound to a shared client.
//
// Example:
//
//	idService := ids.NewService(client)
//	_ = idService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     string
}

// Id returns a handle for a specific id with the given ID.
// This handle can be used to perform operations related to that id,
// such as retrieving profile information and activity data.
//
// Example:
//
//	id := client.Account.Id("uuid-1234")
//	account, err := id.Account(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Level: %d\n", account.Data.Level)
func (s *Service) Id(id string) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type AccountData = experienceClient.AccountResponse
type StreakData = experienceClient.StreakData

// AccountResponse contains the experience account payload for a id.
type AccountResponse struct {
	Data         AccountData
	ResponseMeta common.ResponseMeta
}

// Account retrieves experience and rank information for the given account ID.
//
// Example:
//
//	account, err := client.Account.Id("uuid-1234").Account(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Level: %d (%s)\n", account.Data.Level, account.Data.LevelTitle)
func (h *Handle) Account(ctx context.Context) (AccountResponse, error) {
	resp, err := h.client.ExperienceV1().GetAccountAccountId(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return AccountResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, experienceClient.ParseGetAccountAccountIdResponse)
	if err != nil {
		return AccountResponse{ResponseMeta: meta}, err
	}

	return AccountResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// XPEarnedResponse contains XP history points and HTTP response metadata.
type XPEarnedResponse struct {
	Data         experienceClient.AccountXpEarnedResponse
	ResponseMeta common.ResponseMeta
}

// XPEarned returns experience points earned over the selected period.
//
// Example:
//
//	history, err := client.Account.Id(accountID).XPEarned(ctx, params)
func (h *Handle) XPEarned(ctx context.Context, params *experienceClient.GetAccountXpEarnedParams) (XPEarnedResponse, error) {
	client := h.client
	resp, err := client.ExperienceV1().GetAccountXpEarned(client.Limiter().Wrap(ctx), h.id, params)
	if err != nil {
		return XPEarnedResponse{}, err
	}
	meta := common.ResponseMeta{
		Raw:        extract.Raw(resp),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		CFRay:      resp.Header.Get("CF-Ray"),
	}
	parsed, err := experienceClient.ParseGetAccountXpEarnedResponse(resp)
	if err == nil && meta.StatusCode >= 200 && meta.StatusCode < 300 && parsed.JSON2XX == nil {
		err = errors.New("XP history response is missing its JSON body")
	}
	if err != nil || meta.StatusCode < 200 || meta.StatusCode >= 300 {
		return errutil.UnwrapFailure(err, meta.Raw, meta.StatusCode, func([]byte) XPEarnedResponse {
			return XPEarnedResponse{ResponseMeta: meta}
		})
	}
	return XPEarnedResponse{Data: *parsed.JSON2XX, ResponseMeta: meta}, nil
}
