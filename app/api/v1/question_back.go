package v1

type getQuestionBankVOByIdRequest struct {
	Current     string `json:"current"`
	Description string `json:"description"`
	Id          string `json:"id"`
	NeedQuery   bool   `json:"needQuery"`
	NotId       int64  `json:"notId"`
	PageSize    int32  `json:"pageSize"`
	Picture     string `json:"picture"`
	SearchText  string `json:"searchText"`
	SortField   string `json:"sortField"`
	SortOrder   string `json:"sortOrder"`
	Title       string `json:"title"`
	UserId      int64  `json:"userId"`
}

type getQuestionBankVOByIdResponse struct {
	Response
}
