package types

type Album struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

type Asset struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}
