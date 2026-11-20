package models

type Pokemon struct {
	Name     Names   `json:"name"`
	Id       int     `json:"pokedex_id"`
	Category string  `json:"category"`
	Types    []Types `json:"types"`
	Sprites  Sprites `json:"sprites"`
}

type Names struct {
	Fr string `json:"fr"`
}

type Sprites struct {
	Regular string `json:"regular"`
}

type Types struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
