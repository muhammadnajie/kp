package links

type PayloadCreateLink struct {
	Title   string `json:"title"`
	Address string `json:"address"`
}

type PayloadUpdateLink struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
}

type PayloadDeleteLink struct {
	ID string `json:"id"`
}
