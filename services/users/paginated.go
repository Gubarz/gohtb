package users

import (
	"context"
	"fmt"

	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/ptr"
	"github.com/gubarz/gohtb/internal/service"
)

type UserProfileActivityQuery struct {
	client  service.Client
	page    int
	perPage int
	id      int
}

func (q *UserProfileActivityQuery) Page(n int) *UserProfileActivityQuery {
	qc := ptr.Clone(q)
	qc.page = n
	return qc
}

func (q *UserProfileActivityQuery) PerPage(n int) *UserProfileActivityQuery {
	qc := ptr.Clone(q)
	qc.perPage = n
	return qc
}

func (q *UserProfileActivityQuery) Next() *UserProfileActivityQuery {
	qc := ptr.Clone(q)
	qc.page++
	return qc
}

func (q *UserProfileActivityQuery) Previous() *UserProfileActivityQuery {
	qc := ptr.Clone(q)
	if qc.page > 1 {
		qc.page--
	}
	return qc
}

type UserProfileActivityItems []UserProfileActivity

type UserProfileActivityResponse struct {
	Data         UserProfileActivityItems
	ResponseMeta common.ResponseMeta
}

func (q *UserProfileActivityQuery) Results(ctx context.Context) (UserProfileActivityResponse, error) {
	return q.fetchResults(ctx)
}

func (q *UserProfileActivityQuery) AllResults(ctx context.Context) (UserProfileActivityResponse, error) {
	var all []UserProfileActivity
	page := 1
	var meta common.ResponseMeta

	for {
		qp := ptr.Clone(q)
		qp.page = page

		resp, err := qp.fetchResults(ctx)
		if err != nil {
			return UserProfileActivityResponse{}, err
		}

		all = append(all, resp.Data...)

		meta = resp.ResponseMeta

		if len(resp.Data) < q.perPage {
			break
		}

		page++
	}

	return UserProfileActivityResponse{
		Data:         all,
		ResponseMeta: meta,
	}, nil
}

func (q *UserProfileActivityQuery) fetchResults(ctx context.Context) (UserProfileActivityResponse, error) {
	params := &v5Client.GetUserProfileActivityParams{
		Page:    &q.page,
		PerPage: &q.perPage,
	}

	if q.id == 0 {
		return UserProfileActivityResponse{}, fmt.Errorf("user ID is required")
	}

	resp, err := q.client.V5().GetUserProfileActivity(q.client.Limiter().Wrap(ctx), q.id, params)
	if err != nil {
		return UserProfileActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetUserProfileActivityResponse)
	if err != nil {
		return UserProfileActivityResponse{ResponseMeta: meta}, err
	}

	activities, err := wrapUserProfileActivityItems(parsed.JSON200.Data)
	if err != nil {
		return UserProfileActivityResponse{ResponseMeta: meta}, err
	}

	return UserProfileActivityResponse{
		Data:         activities,
		ResponseMeta: meta,
	}, nil
}

func wrapUserProfileActivityItems(items []v5Client.UserProfileActivityItem) (UserProfileActivityItems, error) {
	out := make(UserProfileActivityItems, len(items))
	for i, item := range items {
		wrapped, err := wrapUserProfileActivityItem(item)
		if err != nil {
			return nil, err
		}
		out[i] = wrapped
	}
	return out, nil
}

func wrapUserProfileActivityItem(item v5Client.UserProfileActivityItem) (UserProfileActivity, error) {
	discriminator, err := item.Discriminator()
	if err != nil {
		return UserProfileActivity{}, err
	}

	switch discriminator {
	case "challenge":
		activity, err := item.AsUserProfileActivityChallenge()
		if err != nil {
			return UserProfileActivity{}, err
		}
		payload := activity
		return UserProfileActivity{
			UserProfileActivityBase: v5Client.UserProfileActivityBase{
				Avatar:  payload.Avatar,
				Blood:   payload.Blood,
				Id:      payload.Id,
				Name:    payload.Name,
				OwnDate: payload.OwnDate,
				Points:  payload.Points,
				Type:    string(payload.Type),
			},
			Challenge: payload,
		}, nil
	case "fortress":
		activity, err := item.AsUserProfileActivityFortress()
		if err != nil {
			return UserProfileActivity{}, err
		}
		payload := activity
		return UserProfileActivity{
			UserProfileActivityBase: v5Client.UserProfileActivityBase{
				Avatar:  payload.Avatar,
				Blood:   payload.Blood,
				Id:      payload.Id,
				Name:    payload.Name,
				OwnDate: payload.OwnDate,
				Points:  payload.Points,
				Type:    string(payload.Type),
			},
			Fortress: payload,
		}, nil
	case "root", "user":
		activity, err := item.AsUserProfileActivityMachineOwn()
		if err != nil {
			return UserProfileActivity{}, err
		}
		payload := activity
		return UserProfileActivity{
			UserProfileActivityBase: v5Client.UserProfileActivityBase{
				Avatar:  payload.Avatar,
				Blood:   payload.Blood,
				Id:      payload.Id,
				Name:    payload.Name,
				OwnDate: payload.OwnDate,
				Points:  payload.Points,
				Type:    string(payload.Type),
			},
			MachineOwn: payload,
		}, nil
	case "prolab":
		activity, err := item.AsUserProfileActivityProlab()
		if err != nil {
			return UserProfileActivity{}, err
		}
		payload := activity
		return UserProfileActivity{
			UserProfileActivityBase: v5Client.UserProfileActivityBase{
				Avatar:  payload.Avatar,
				Blood:   payload.Blood,
				Id:      payload.Id,
				Name:    payload.Name,
				OwnDate: payload.OwnDate,
				Points:  payload.Points,
				Type:    string(payload.Type),
			},
			Prolab: payload,
		}, nil
	case "sherlock":
		activity, err := item.AsUserProfileActivitySherlock()
		if err != nil {
			return UserProfileActivity{}, err
		}
		payload := activity
		return UserProfileActivity{
			UserProfileActivityBase: v5Client.UserProfileActivityBase{
				Avatar:  payload.Avatar,
				Blood:   payload.Blood,
				Id:      payload.Id,
				Name:    payload.Name,
				OwnDate: payload.OwnDate,
				Points:  payload.Points,
				Type:    string(payload.Type),
			},
			Sherlock: payload,
		}, nil
	default:
		return UserProfileActivity{}, fmt.Errorf("unsupported user profile activity type %q", discriminator)
	}
}
