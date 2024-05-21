package components

type Pagination struct {
	PageStartRecord int
	PageEndRecord   int
	TotalRecords    int
	PrevLink        string
	NextLink        string
}
