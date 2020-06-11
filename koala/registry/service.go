package registry


type Service struct {
	Name string `json:"name"`
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	Id string `json:"id"`
	IP string `json:"ip"`
	Port int `json:"port"`
	Weight int `json:"weight"`
}