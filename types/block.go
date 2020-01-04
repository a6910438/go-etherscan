package types

type Block struct {
	ID     int `storm:"id,increment"`
	Height int
}
