package sherlocks

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromAPISherlockList(data v4client.SherlockItem) SherlockItem {
	return SherlockItem{
		AuthUserHasReviewed: deref.Bool(data.AuthUserHasReviewed),
		Avatar:              deref.String(data.Avatar),
		CategoryId:          deref.Int(data.CategoryId),
		CategoryName:        deref.String(data.CategoryName),
		Difficulty:          deref.String(data.Difficulty),
		Id:                  deref.Int(data.Id),
		IsOwned:             deref.Bool(data.IsOwned),
		Name:                deref.String(data.Name),
		Pinned:              deref.Bool(data.Pinned),
		PlayMethods:         fromStringArray(data.PlayMethods),
		Progress:            deref.Int(data.Progress),
		Rating:              deref.Float32(data.Rating),
		RatingCount:         deref.Int(data.RatingCount),
		ReleaseDate:         deref.String(data.ReleaseDate),
		Retires:             deref.Time(&data.Retires.Time),
		Solves:              deref.Int(data.Solves),
		State:               deref.String(data.State),
	}
}

func fromAPISherlock(data *v4client.SherlockNamedItemData) SherlockNamedItemData {
	if data == nil {
		return SherlockNamedItemData{}
	}
	return SherlockNamedItemData{
		AuthUserHasReviewed: deref.Bool(data.AuthUserHasReviewed),
		Avatar:              deref.String(data.Avatar),
		CategoryId:          deref.Int(data.CategoryId),
		CategoryName:        deref.String(data.CategoryName),
		Difficulty:          deref.String(data.Difficulty),
		Favorite:            deref.Bool(data.Favorite),
		Id:                  deref.Int(data.Id),
		IsTodo:              deref.Bool(data.IsTodo),
		Name:                deref.String(data.Name),
		PlayMethods:         fromStringArray(data.PlayMethods),
		Rating:              deref.Float32(data.Rating),
		RatingCount:         deref.Int(data.RatingCount),
		ReleaseAt:           deref.Time(data.ReleaseAt),
		Retired:             deref.Bool(data.Retired),
		ShowGoVip:           deref.Bool(data.ShowGoVip),
		State:               deref.String(data.State),
		Tags:                fromStringArray(data.Tags),
		UserCanReview:       deref.Bool(data.UserCanReview),
		WriteupVisible:      deref.Bool(data.WriteupVisible),
	}
}

func fromStringArray(data *v4client.StringArray) []string {
	if data == nil {
		return nil
	}
	return []string(*data)
}
