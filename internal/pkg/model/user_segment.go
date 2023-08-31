package model

type UserSegmentRequest struct {
	UserID                 *int      `json:"user_id,omitempty"`
	SegmentsTitlesToAdd    *[]string `json:"seg_titles_add,omitempty"`
	SegmentsTitlesToDelete *[]string `json:"seg_titles_delete,omitempty"`
	TTL                    *[]int    `json:"ttl,omitempty"`
}

type UserSegmentResponse struct {
	Status                       string         `json:"status"`
	Segments                     []UserSegments `json:"segments,omitempty"`
	SegmentsTitlesNotExistAdd    []string       `json:"seg_titles_not_exist_add,omitempty"`
	SegmentsTitlesNotExistDelete []string       `json:"seg_titles_not_exist_delete,omitempty"`
}

type UserSegments struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
