package swagger

type MethodProperties struct {
	Responses *Code_200 `json:"responses"`
	ResourceDescription string `json:"description"`
	Produces []string `json:"produces"`
	Consumes []string `json:"consumes"`
	Summary string `json:"summary"`
}

