package fortresses

import (
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
)

func fromFortressData(data *v4client.FortressData) Data {
	return Data{
		Company:              fromFortressCompanyData(data.Company),
		CompletionMessage:    deref.String(data.CompletionMessage),
		CoverImageUrl:        deref.String(data.CoverImageUrl),
		Description:          deref.String(data.Description),
		Flags:                convert.Slice(*data.Flags, common.FromAPIFlag),
		HasCompletionMessage: deref.Bool(data.HasCompletionMessage),
		Id:                   deref.Int(data.Id),
		Image:                deref.String(data.Image),
		Ip:                   deref.String(data.Ip),
		Name:                 deref.String(data.Name),
		PlayersCompleted:     deref.Int(data.PlayersCompleted),
		Points:               deref.String(data.Points),
		ProgressPercent:      deref.Float32(data.ProgressPercent),
		ResetVotes:           deref.Int(data.ResetVotes),
		UserAvailability:     fromFortressUserAvailabilityData(data.UserAvailability),
	}
}

func fromFortressCompanyData(data *v4client.Company) Company {
	return Company{
		Description: deref.String(data.Description),
		Id:          deref.Int(data.Id),
		Image:       deref.String(data.Image),
		Name:        deref.String(data.Name),
		Url:         deref.String(data.Url),
	}
}

func fromFortressUserAvailabilityData(data *v4client.UserAvailability) UserAvailability {
	return UserAvailability{
		Available: deref.Bool(data.Available),
		Code:      deref.Int(data.Code),
		Message:   deref.String(data.Message),
	}
}

func toItem(data v4client.Fortress) Item {
	return Item{
		CoverImageUrl: deref.String(data.CoverImageUrl),
		Id:            deref.Int(data.Id),
		Image:         deref.String(data.Image),
		Name:          deref.String(data.Name),
		New:           deref.Bool(data.New),
		NumberOfFlags: deref.Int(data.NumberOfFlags),
	}
}
