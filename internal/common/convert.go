package common

import (
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
)

func FromAPIUserBasicInfo(data v4client.UserBasicInfo) UserBasicInfo {
	return UserBasicInfo{
		Avatar: deref.String(data.Avatar),
		Id:     deref.Int(data.Id),
		Name:   deref.String(data.Name),
	}
}

func FromAPILinks(data v4client.Links) Links {
	return Links{
		First: deref.String(data.First),
		Last:  deref.String(data.Last),
		Next:  deref.String(data.Next),
		Prev:  deref.String(data.Prev),
	}
}

func FromAPIMeta(data v4client.Meta) Meta {
	return Meta{
		CurrentPage: deref.Int(data.CurrentPage),
		From:        deref.Int(data.From),
		LastPage:    deref.Int(data.LastPage),
		Links:       convert.Slice(*data.Links, fromAPIPaginationLink),
		Path:        deref.String(data.Path),
		PerPage:     deref.Int(data.PerPage),
		To:          deref.Int(data.To),
		Total:       deref.Int(data.Total),
	}
}

func FromAPIMetaAlt(data v4client.MetaAlt) MetaAlt {
	return MetaAlt{
		CurrentPage: deref.Int(data.CurrentPage),
		Pages:       deref.Int(data.Pages),
	}
}

func fromAPIPaginationLink(data v4client.PaginationLink) PaginationLink {
	return PaginationLink{
		Active: deref.Bool(data.Active),
		Label:  deref.String(data.Label),
		Url:    deref.String(data.Url),
	}
}

func fromAPITag(data v4client.Tag) Tag {
	return Tag{
		Id:            deref.Int(data.Id),
		Name:          deref.String(data.Name),
		TagCategoryId: deref.Int(data.TagCategoryId),
	}
}

func FromAPITagCategory(data v4client.TagCategory) TagCategory {
	return TagCategory{
		Id:   deref.Int(data.Id),
		Name: deref.String(data.Name),
		Tags: convert.Slice(*data.Tags, fromAPITag),
	}
}

func FromAPIHelpfulReviews(data v4client.HelpfulReviews) HelpfulReviews {
	return HelpfulReviews{
		Id:       deref.Int(data.Id),
		ReviewId: deref.Int(data.ReviewId),
		UserId:   deref.Int(data.UserId),
	}
}

func FromAPIDifficultyChart(data *v4client.DifficultyChart) DifficultyChart {
	values, err := data.AsDifficultyChart1()
	if err != nil {
		return DifficultyChart{
			CounterBitHard:   0,
			CounterBrainFuck: 0,
			CounterCake:      0,
			CounterEasy:      0,
			CounterExHard:    0,
			CounterHard:      0,
			CounterMedium:    0,
			CounterTooEasy:   0,
			CounterTooHard:   0,
			CounterVeryEasy:  0,
		}
	}
	return DifficultyChart{
		CounterBitHard:   deref.Int(values.CounterBitHard),
		CounterBrainFuck: deref.Int(values.CounterBrainFuck),
		CounterCake:      deref.Int(values.CounterCake),
		CounterEasy:      deref.Int(values.CounterEasy),
		CounterExHard:    deref.Int(values.CounterExHard),
		CounterHard:      deref.Int(values.CounterHard),
		CounterMedium:    deref.Int(values.CounterMedium),
		CounterTooEasy:   deref.Int(values.CounterTooEasy),
		CounterTooHard:   deref.Int(values.CounterTooHard),
		CounterVeryEasy:  deref.Int(values.CounterVeryEasy),
	}
}

func FromAPIInfoArray(data v4client.Item) Item {
	return Item{
		Id: deref.Int(data.Id),
	}
}

func FromAPIPlayInfoAlt(data v4client.PlayInfoAlt) PlayInfoAlt {
	return PlayInfoAlt{
		ExpiresAt: deref.Time(data.ExpiresAt),
		Ip:        deref.String(data.Ip),
		Ports:     deref.Slice(data.Ports),
		Status:    deref.String(data.Status),
	}
}

func FromAPIBloodInfo(data *v4client.BloodInfo) BloodInfo {
	if data == nil {
		return BloodInfo{
			BloodDifference: "",
			CreatedAt:       "",
			User:            UserBasicInfo{},
		}
	}
	return BloodInfo{
		BloodDifference: deref.String(data.BloodDifference),
		CreatedAt:       deref.String(data.CreatedAt),
		User:            FromAPIUserBasicInfo(*data.User),
	}
}

func FromAPIMaker(data *v4client.Maker) Maker {
	if data == nil {
		return Maker{
			Avatar:      "",
			Id:          0,
			IsRespected: false,
			Name:        "",
		}
	}
	return Maker{
		Avatar:      deref.String(data.Avatar),
		Id:          deref.Int(data.Id),
		IsRespected: deref.Bool(data.IsRespected),
		Name:        deref.String(data.Name),
	}
}

func FromAPIMatrixInfo(data v4client.MatrixInfo) MatrixInfo {
	return MatrixInfo{
		Ctf:    deref.Float32(data.Ctf),
		Custom: deref.Float32(data.Custom),
		Cve:    deref.Float32(data.Cve),
		Enum:   deref.Float32(data.Enum),
		Real:   deref.Float32(data.Real),
	}
}

func fromAPIAcademyDifficulty(data v4client.AcademyDifficulty) AcademyDifficulty {
	return AcademyDifficulty{
		Color: deref.String(data.Color),
		Id:    deref.Int(data.Id),
		Level: deref.Int(data.Level),
		Text:  deref.String(data.Text),
		Title: deref.String(data.Title),
	}
}

func FromAPIAcademyModule(data v4client.AcademyModule) AcademyModule {
	return AcademyModule{
		Aggregates: deref.Slice(data.Aggregates),
		Avatar:     deref.String(data.Avatar),
		Difficulty: fromAPIAcademyDifficulty(*data.Difficulty),
		Id:         deref.Int(data.Id),
		Logo:       deref.String(data.Logo),
		Name:       deref.String(data.Name),
		Tier:       fromAPIAcademyTiers(*data.Tier),
		Url:        deref.String(data.Url),
	}
}

func fromAPIAcademyTiers(data v4client.AcademyTiers) AcademyTiers {
	return AcademyTiers{
		Color:  deref.String(data.Color),
		Name:   deref.String(data.Name),
		Number: deref.Int(data.Number),
	}
}

func FromAPITeamMachineAttackPaths(data v4client.TeamMachineAttackPaths) TeamMachineAttackPaths {
	return TeamMachineAttackPaths{
		BinaryAnalysis:            fromAPITeamsAttackPathCard(*data.BinaryAnalysis),
		BinaryExploitation:        fromAPITeamsAttackPathCard(*data.BinaryExploitation),
		ConfigurationAnalysis:     fromAPITeamsAttackPathCard(*data.ConfigurationAnalysis),
		Fuzzing:                   fromAPITeamsAttackPathCard(*data.Fuzzing),
		Impersonation:             fromAPITeamsAttackPathCard(*data.Impersonation),
		PacketCaptureAnalysis:     fromAPITeamsAttackPathCard(*data.PacketCaptureAnalysis),
		Pivoting:                  fromAPITeamsAttackPathCard(*data.Pivoting),
		Reconnaissance:            fromAPITeamsAttackPathCard(*data.Reconnaissance),
		UserEnumeration:           fromAPITeamsAttackPathCard(*data.UserEnumeration),
		WebSiteStructureDiscovery: fromAPITeamsAttackPathCard(*data.WebSiteStructureDiscovery),
	}
}

func fromAPITeamsAttackPathCard(data v4client.TeamsAttackPathCard) TeamsAttackPathCard {
	return TeamsAttackPathCard{
		AvgTeamsSolved: deref.Float32(data.AvgTeamsSolved),
		Name:           deref.String(data.Name),
		Solved:         deref.Int(data.Solved),
		Total:          deref.Int(data.Total),
	}
}
