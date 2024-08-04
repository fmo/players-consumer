package models

type Player struct {
	Id                  string `json:"id"`
	Team                string `json:"team"`
	TeamId              int32  `json:"teamId"`
	Name                string `json:"name"`
	Firstname           string `json:"firstname"`
	Lastname            string `json:"lastname"`
	Age                 int32  `json:"age"`
	Nationality         string `json:"nationality"`
	Photo               string `json:"photo"`
	ApiFootballId       string `json:"apiFootballId"`
	Appearances         int32  `json:"appearances"`
	Position            string `json:"position"`
	TransfermarktId     string `json:"transfermarktId"`
	ShirtNumber         string `json:"shirtNumber"`
	MarketValue         int32  `json:"marketValue"`
	MarketValueCurrency string `json:"marketValueCurrency"`
}
