package platform

type Fightcade struct{}

func (s Fightcade) Search(query string) ([]SearchResult, error) {
	results := make([]SearchResult, 2)
	results[0] = SearchResult{Game: "Street Fighter 2: Champion Edition", AppID: "sf2ce"}
	results[1] = SearchResult{Game: "Garou: Mark of the Wolves", AppID: "garou"}
	return results, nil
}

func (s Fightcade) Launch(appId string) error {
	//TODO implement me
	panic("implement me")
}
