package data

type Classification struct{
	Name string `json:"name" cql:"name"`
	ProbaHateful float64 `json:"proba_hateful" cql:"proba_hateful"`
	ProbaNotHateful float64 `json:"proba_not_hateful" cql:"proba_not_hateful"`
	Class int `json:"class" cql:"class"`
}
