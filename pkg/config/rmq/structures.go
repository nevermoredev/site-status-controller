package rmq




type JobsRequest struct {
	PageId string `json:PageId`
	Url    string `json:Url`
}