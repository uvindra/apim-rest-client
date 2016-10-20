package swagger

type InfoType struct {
	Title string `json:"title"`
	Version string `json:"version"`
}

type Swagger struct {
	Version string `json:"swagger"`
	Paths *PathsType `json:"paths"`
	Info *InfoType `json:"info"`
}
