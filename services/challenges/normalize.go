package challenges

import (
	"strconv"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromAPIChallengeList(data v4client.ChallengeList) ChallengeList {
	return ChallengeList{
		AuthUserHasReviewed: deref.Bool(data.AuthUserHasReviewed),
		Avatar:              deref.String(data.Avatar),
		CategoryId:          deref.Int(data.CategoryId),
		CategoryName:        deref.String(data.CategoryName),
		Difficulty:          deref.String(data.Difficulty),
		DifficultyChart:     common.FromAPIDifficultyChart(data.DifficultyChart),
		Id:                  deref.Int(data.Id),
		IsOwned:             deref.Bool(data.IsOwned),
		Name:                deref.String(data.Name),
		Pinned:              deref.Bool(data.Pinned),
		PlayMethods:         deref.Slice(data.PlayMethods),
		Rating:              deref.Float32(data.Rating),
		RatingCount:         deref.Int(data.RatingCount),
		ReleaseDate:         deref.Time(data.ReleaseDate),
		Retires:             fromAPIChallengeRetires(data.Retires),
		Solves:              deref.Int(data.Solves),
		State:               deref.String(data.State),
		UserDifficulty:      deref.String(data.UserDifficulty),
	}
}

func fromAPIChallengeRetires(data *v4client.ChallengeRetires) ChallengeRetires {
	if data == nil {
		return ChallengeRetires{}
	}
	return ChallengeRetires{
		name:       deref.String(data.Name),
		difficulty: deref.String(data.Difficulty),
	}
}

func fromAPIChallengeActivity(data v4client.ChallengeActivity) ChallengeActivity {
	return ChallengeActivity{
		CreatedAt:  deref.Time(data.CreatedAt),
		Date:       deref.String(data.Date),
		DateDiff:   deref.String(data.DateDiff),
		Type:       deref.String(data.Type),
		UserAvatar: deref.String(data.UserAvatar),
		UserId:     deref.Int(data.UserId),
		UserName:   deref.String(data.UserName),
	}
}

func fromAPIChallengeInfo(data *v4client.Challenge) Challenge {
	if data == nil {
		return Challenge{}
	}
	return Challenge{
		AuthUserHasReviewed:  deref.Bool(data.AuthUserHasReviewed),
		AuthUserSolve:        deref.Bool(data.AuthUserSolve),
		AuthUserSolveTime:    deref.String(data.AuthUserSolveTime),
		CanAccessWalkthough:  deref.Bool(data.CanAccessWalkthough),
		CategoryName:         deref.String(data.CategoryName),
		Creator2Avatar:       deref.String(data.Creator2Avatar),
		Creator2Id:           deref.Int(data.Creator2Id),
		Creator2Name:         deref.String(data.Creator2Name),
		CreatorAvatar:        deref.String(data.CreatorAvatar),
		CreatorId:            deref.Int(data.CreatorId),
		CreatorName:          deref.String(data.CreatorName),
		Description:          deref.String(data.Description),
		Difficulty:           deref.String(data.Difficulty),
		DifficultyChart:      common.FromAPIDifficultyChart(data.DifficultyChart),
		DislikeByAuthUser:    deref.Bool(data.DislikeByAuthUser),
		Dislikes:             deref.Int(data.Dislikes),
		Docker:               deref.Bool(data.Docker),
		DockerIp:             deref.String(data.DockerIp),
		DockerPorts:          deref.String(data.DockerPorts),
		DockerStatus:         deref.String(data.DockerStatus),
		Download:             deref.Bool(data.Download),
		FirstBloodTime:       deref.String(data.FirstBloodTime),
		FirstBloodUser:       deref.String(data.FirstBloodUser),
		FirstBloodUserAvatar: deref.String(data.FirstBloodUserAvatar),
		FirstBloodUserId:     deref.Int(data.FirstBloodUserId),
		FileName:             deref.String(data.FileName),
		FileSize:             deref.String(data.FileSize),
		HasChangelog:         deref.Bool(data.HasChangelog),
		Id:                   deref.Int(data.Id),
		IsRespected:          deref.Bool(data.IsRespected),
		IsRespected2:         deref.Bool(data.IsRespected2),
		IsTodo:               deref.Bool(data.IsTodo),
		LikeByAuthUser:       deref.Bool(data.LikeByAuthUser),
		Likes:                deref.Int(data.Likes),
		Name:                 deref.String(data.Name),
		PlayInfo:             common.FromAPIPlayInfoAlt(data.PlayInfo),
		PlayMethods:          deref.Slice(data.PlayMethods),
		Points:               pointsToInt(data.Points),
		Recommended:          deref.Int(data.Recommended),
		ReleaseDate:          deref.Time(data.ReleaseDate),
		Released:             deref.Int(data.Released),
		Retired:              deref.Bool(data.Retired),
		ReviewsCount:         deref.Int(data.ReviewsCount),
		Sha256:               deref.String(data.Sha256),
		ShowGoVip:            deref.Bool(data.ShowGoVip),
		Solves:               deref.Int(data.Solves),
		Stars:                deref.Float32(data.Stars),
		State:                deref.String(data.State),
		Tags:                 deref.Slice(data.Tags),
		UserCanReview:        deref.Bool(data.UserCanReview),
	}
}

func pointsToInt(v *v4client.Challenge_Points) int {
	if v == nil {
		return 0
	}

	if b, err := v.AsChallengePoints0(); err == nil {
		return b
	}
	if b, err := v.AsChallengePoints1(); err == nil {
		if v, err := strconv.Atoi(b); err == nil {
			return v
		}
		return 0
	}
	return 0
}
