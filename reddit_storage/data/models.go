package data

type Classification struct {
	Name            string  `json:"name" cql:"name"`
	Type            string  `json:"type" cql:"type"`
	ProbaHateful    float64 `json:"proba_hateful"`
	ProbaNotHateful float64 `json:"proba_not_hateful"`
	Class           int     `json:"is_hatespech"`
}
