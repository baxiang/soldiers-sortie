package utils


type Pagination struct {
	Count     int
	PageIndex int
	PageSize  int
	Sorter    []string
	Data      interface{}
}

// Image struct
type Image struct {
	Filepath []byte
	Md5      string
}

// Images muti imgs.
type Images []*Image