package dto

type Filter struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Deleted bool `json:"deleted"`
}
