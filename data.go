package trello

//Auth type to hold authentication data for trello API
type Auth struct {
	Key   string
	Token string
}

type trelloData struct {
	List struct {
		Id   string
		Name string
	} `json:"list"`
	Board struct {
		Id        string
		Name      string
		ShortLink string
	} `json:"board"`
	Card struct {
		Id        string
		Name      string
		ShortLink string
		IdShort   int
	} `json:"card"`
	Text string `json:"text"`
}

type trelloAction struct {
	Id            string
	Data          *trelloData
	MemberCreator *trelloMemberCreator
}

type trelloMemberCreator struct {
	Id       string
	FullName string
	Username string
}

type trelloCard struct {
	IdF               string `json:"id"`
	DateLastActivityF string `json:"dateLastActivity"`
	DescF             string `json:"desc"`
	IdBoardF          string `json:"idBoard"`
	IdListF           string `json:"idList"`
	NameF             string `json:"name"`
	ShortUrlF         string `json:"shortUrl"`
}

//Comment represents a minimal Trello Action of type comment
type Comment interface {
	Text() string
	Author() string
}

func (a trelloAction) Author() string {
	if a.MemberCreator == nil {
		return ""
	} else {
		return a.MemberCreator.FullName
	}
}

func (a trelloAction) Text() string {
	if a.Data == nil {
		return ""
	} else {
		return a.Data.Text
	}
}

//Card represents a minimal Trello Card
type Card interface {
	Id() string
	DateLastActivity() string
	IdBoard() string
	IdList() string
	Name() string
	Desc() string
	ShortUrl() string
}

func (i trelloCard) Id() string {
	return i.IdF
}

func (i trelloCard) IdList() string {
	return i.IdListF
}

func (i trelloCard) Name() string {
	return i.NameF
}

func (i trelloCard) Desc() string {
	return i.DescF
}

func (i trelloCard) IdBoard() string {
	return i.IdBoardF
}
func (i trelloCard) DateLastActivity() string {
	return i.DateLastActivityF
}

func (i trelloCard) ShortUrl() string {
	return i.ShortUrlF
}

//Board represents a minimal Trello Board
type Board interface {
	Name() string
	Id() string
}

type trelloBoard struct {
	IdF   string `json:"id"`
	NameF string `json:"name"`
}

func (i trelloBoard) Id() string {
	return i.IdF
}

func (i trelloBoard) Name() string {
	return i.NameF
}

//List represents a minimal Trello List
type List interface {
	Name() string
	Id() string
	Closed() bool
}

type trelloList struct {
	IdF         string `json:"id"`
	NameF       string `json:"name"`
	ClosedF     bool   `json:"closed"`
	IdBoardF    string `json:"idBoard"`
	SubscribedF bool   `json:"subscribed"`
	PosF        int    `json:"pos"`
}

func (i trelloList) Id() string {
	return i.IdF
}

func (i trelloList) Name() string {
	return i.NameF
}

func (i trelloList) Closed() bool {
	return i.ClosedF
}
