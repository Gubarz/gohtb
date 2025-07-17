package prolabs

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromAPIProlab(data v4client.Prolab) Prolab {
	return Prolab{
		CoverImgUrl:                deref.String(data.CoverImgUrl),
		DesignatedCategory:         deref.String(data.DesignatedCategory),
		Id:                         deref.Int(data.Id),
		Identifier:                 deref.String(data.Identifier),
		LabServersCount:            deref.Int(data.LabServersCount),
		Level:                      deref.Int(data.Level),
		Mini:                       deref.Bool(data.Mini),
		Name:                       deref.String(data.Name),
		New:                        deref.Bool(data.New),
		Ownership:                  deref.Float32(data.Ownership),
		ProFlagsCount:              deref.Int(data.ProFlagsCount),
		ProMachinesCount:           deref.Int(data.ProMachinesCount),
		ReleaseAt:                  deref.Time(data.ReleaseAt),
		SkillLevel:                 deref.String(data.SkillLevel),
		State:                      deref.String(data.State),
		Team:                       deref.String(data.Team),
		UserEligibleForCertificate: deref.Bool(data.UserEligibleForCertificate),
	}
}

func fromAPIFaqItem(data v4client.FaqItem) FaqItem {
	return FaqItem{
		Generic:  deref.Bool(data.Generic),
		Question: deref.String(data.Question),
		Answer:   deref.String(data.Answer),
	}
}

func fromAPIProlabData(data *v4client.ProlabData) ProlabData {
	if data == nil {
		return ProlabData{}
	}
	return ProlabData{
		ActiveUsers:      deref.Int(data.ActiveUsers),
		CanInteract:      deref.Bool(data.CanInteract),
		CoverImageUrl:    deref.String(data.CoverImageUrl),
		Description:      deref.String(data.Description),
		EntryPoints:      deref.Slice(data.EntryPoints),
		Id:               deref.Int(data.Id),
		Identifier:       deref.String(data.Identifier),
		LabMasters:       convert.SlicePointer(data.LabMasters, common.FromAPIUserIdNameThumb),
		LabServersCount:  deref.Int(data.LabServersCount),
		Mini:             deref.Bool(data.Mini),
		Name:             deref.String(data.Name),
		ProFlagsCount:    deref.Int(data.ProFlagsCount),
		ProMachinesCount: deref.Int(data.ProMachinesCount),
		State:            deref.String(data.State),
		Version:          deref.String(data.Version),
		VideoUrl:         deref.String(data.VideoUrl),
		Writeup:          deref.String(data.Writeup),
	}
}

func fromAPIProlabMachineData(data v4client.Machine) Machine {
	return Machine{
		AvatarThumbUrl: deref.String(data.AvatarThumbUrl),
		Id:             deref.Int(data.Id),
		Name:           deref.String(data.Name),
		Os:             deref.String(data.Os),
	}
}

func fromAPIProlabOverviewData(data *v4client.ProlabOverviewData) ProlabOverviewData {
	if data == nil {
		return ProlabOverviewData{}
	}
	return ProlabOverviewData{
		DesignatedLevel:    fromAPIDesignatedLevel(data.DesignatedLevel),
		Excerpt:            deref.String(data.Excerpt),
		Id:                 deref.Int(data.Id),
		Identifier:         deref.String(data.Identifier),
		LabMasters:         convert.SlicePointer(data.LabMasters, common.FromAPIUserIdNameThumb),
		Mini:               deref.Bool(data.Mini),
		Name:               deref.String(data.Name),
		NewVersion:         deref.Bool(data.NewVersion),
		OverviewImageUrl:   deref.String(data.OverviewImageUrl),
		ProFlagsCount:      deref.Int(data.ProFlagsCount),
		ProMachinesCount:   deref.Int(data.ProMachinesCount),
		SkillLevel:         deref.String(data.SkillLevel),
		SocialLinks:        fromAPISocialLinks(data.SocialLinks),
		State:              deref.String(data.State),
		UserEligibleToPlay: deref.Bool(data.UserEligibleToPlay),
		Version:            deref.String(data.Version),
	}
}

func fromAPIDesignatedLevel(data *v4client.DesignatedLevel) DesignatedLevel {
	if data == nil {
		return DesignatedLevel{}
	}
	return DesignatedLevel{
		Category:    deref.String(data.Category),
		Description: deref.String(data.Description),
		Level:       deref.Int(data.Level),
		Team:        deref.String(data.Team),
	}
}

func fromAPISocialLinks(data *v4client.SocialLinks) SocialLinks {
	if data == nil {
		return SocialLinks{}
	}
	return SocialLinks{
		Discord: deref.String(data.Discord),
		Forum:   deref.String(data.Forum),
	}
}

func fromAPIProlabProgressData(data *v4client.ProlabProgressData) ProlabProgressData {
	if data == nil {
		return ProlabProgressData{}
	}
	return ProlabProgressData{
		KeyedProLabMileStone:              convert.SlicePointer(data.KeyedProLabMileStone, fromAPIKeyedProLabMileStone),
		Ownership:                         deref.Float32(data.Ownership),
		OwnershipRequiredForCertification: deref.Float32(data.OwnershipRequiredForCertification),
		UserEligibleForCertificate:        deref.Bool(data.UserEligibleForCertificate),
	}
}

func fromAPIKeyedProLabMileStone(data v4client.KeyedProLabMileStone) KeyedProLabMileStone {
	return KeyedProLabMileStone{
		Description:        deref.String(data.Description),
		Icon:               deref.String(data.Icon),
		IsMilestoneReached: deref.Bool(data.IsMilestoneReached),
		Percent:            deref.Float32(data.Percent),
		Rarity:             deref.Float32(data.Rarity),
		Text:               deref.String(data.Text),
	}
}

func fromAPIProlabSubscription(data *v4client.ProlabSubscription) ProlabSubscription {
	if data == nil {
		return ProlabSubscription{}
	}
	return ProlabSubscription{
		Active:             deref.Bool(data.Active),
		EndsAt:             deref.String(data.EndsAt),
		Name:               deref.String(data.Name),
		RenewsAt:           deref.String(data.RenewsAt),
		SubscriptionPeriod: deref.String(data.SubscriptionPeriod),
		Type:               deref.String(data.Type),
	}
}
