package todo

type Filters struct {
	Status string `query:"status"`
	Q      string `query:"q"`
}
