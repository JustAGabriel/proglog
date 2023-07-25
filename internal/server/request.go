package server

type CreateRecordRequest struct {
	Record Record `json:"record"`
}

type CreateRecordResponse struct {
	Offset uint `json:"offset"`
}

type GetRecordRequest struct {
	Offset uint `json:"offset"`
}

type GetRecordResponse struct {
	Record Record `json:"record"`
}
