package dtos

type CommentResp struct {
	List     []TComment `json:"lists"`
	PageInfo PageInfo   `json:"pagination"`
}
