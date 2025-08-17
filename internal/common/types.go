package common

import (
	"net/http"
	"time"
)

type MessageResponse struct {
	Data Message
	ResponseMeta
}

type ResponseMeta struct {
	Raw        []byte
	StatusCode int
	Headers    http.Header
	CFRay      string
}

type FlagData struct {
	Flags  []Flag
	Status bool
	ResponseMeta
}
type InfoArray = []Item

type Item struct {
	Id int
}

type TodoUpdateResponse struct {
	Data InfoArray
	ResponseMeta
}

type Flag struct {
	Id     int
	Owned  bool
	Points int
	Title  string
}

type TodoItem struct {
	Id int
}

type Message struct {
	Message string
	Success bool
}

type AcademyDifficulty struct {
	Color string
	Id    int
	Level int
	Text  string
	Title string
}

type AcademyModule struct {
	Avatar     string
	Difficulty AcademyDifficulty
	Id         int
	Logo       string
	Name       string
	Tier       AcademyTiers
	Url        string
}

type AcademyTiers struct {
	Color  string
	Name   string
	Number int
}

type DifficultyChart struct {
	CounterBitHard   int
	CounterBrainFuck int
	CounterCake      int
	CounterEasy      int
	CounterExHard    int
	CounterHard      int
	CounterMedium    int
	CounterTooEasy   int
	CounterTooHard   int
	CounterVeryEasy  int
}

type Links struct {
	First string
	Last  string
	Next  string
	Prev  string
}

type Maker struct {
	Avatar      string
	Id          int
	IsRespected bool
	Name        string
}

type MatrixInfo struct {
	Ctf    float32
	Custom float32
	Cve    float32
	Enum   float32
	Real   float32
}

type Meta struct {
	CurrentPage int
	From        int
	LastPage    int
	Links       []PaginationLink
	Path        string
	PerPage     int
	To          int
	Total       int
}

type MetaAlt struct {
	CurrentPage int
	Pages       int
}

type PaginationLink struct {
	Active bool
	Label  string
	Url    string
}

type UserBasicInfo struct {
	Avatar string
	Id     int
	Name   string
}

type UserBasicInfoWithRespect struct {
	Avatar      string
	Id          int
	Name        string
	IsRespected bool
}

type BloodInfo struct {
	BloodDifference string
	CreatedAt       string
	User            UserBasicInfo
}

type TeamMachineAttackPaths struct {
	Blockchain              TeamsAttackPathCard
	Cloud                   TeamsAttackPathCard
	EnterpriseNetwork       TeamsAttackPathCard
	Forensics               TeamsAttackPathCard
	Mobile                  TeamsAttackPathCard
	NicheTechnologies       TeamsAttackPathCard
	Person                  TeamsAttackPathCard
	SecurityOperations      TeamsAttackPathCard
	VulnerabilityAssessment TeamsAttackPathCard
	WebApplication          TeamsAttackPathCard
}

type TeamsAttackPathCard struct {
	AvgTeamsSolved float32
	Name           string
	Solved         int
	Total          int
}

type HelpfulReviews struct {
	Id       int
	ReviewId int
	UserId   int
}

type Tag struct {
	Id            int
	Name          string
	TagCategoryId int
}

type TagCategory struct {
	Id   int
	Name string
	Tags []Tag
}

type PlayInfo struct {
	ActivePlayerCount int
	ExpiresAt         time.Time
	IsActive          bool
	IsSpawned         bool
	IsSpawning        bool
}

type OwnsResponse struct {
	Data string
	ResponseMeta
}

type Messagesuccess struct {
	Message string
	Success bool
}

type PlayInfoAlt struct {
	ExpiresAt time.Time
	Ip        string
	Ports     IntArray
	Status    string
}

type IntArray = []int
type StringArray = []string

type UserIdNameThumb struct {
	AvatarThumb string
	Id          int
	Name        string
}
