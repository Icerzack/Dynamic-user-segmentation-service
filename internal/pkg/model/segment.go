package model

type SegmentRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type SegmentResponse struct {
	Status string `json:"status"`
}
