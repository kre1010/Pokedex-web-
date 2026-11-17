package models	


type Pokemon struct {
  Name Names   `json:"name"`
   Id int      `json:"pokemon_id"`
  Category string `json:"category"`
  Types []Types `json:"types"`
  Image Sprites `json:"image"`
}
  
type Names struct {
    Fr string `json:"fr"`
}

type Sprites struct {
    Regular string `json:"regular"`
}

type Types struct {
     Name string `json:"regular"`
	 Image string `json:"image"`
}