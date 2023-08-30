package dto

// UserDTO godoc
//
// swagger:model
type UserDTO struct {
	// UserID
	// required: true
	// min: 1
	// example: 75
	UserID int `json:"user_id"`
}

// UserSegmentsInputDTO godoc
//
// swagger:model
type UserSegmentsInputDTO struct {
	// UserID
	// required: true
	// min: 1
	// example: 75
	UserID int `json:"user_id"`

	// OldSegmentNames - segment names to removing
	// required: false
	// min items: 0
	// example: ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]
	OldSegmentNames []string `json:"old_segment_names,omitempty"`

	// NewSegmentNames - segment names to adding
	// required: false
	// min items: 0
	// example: ["AVITO_VOICE_MESSAGES", "AVITO_DISCOUNT_50"]
	NewSegmentNames []string `json:"new_segment_names,omitempty"`

	// TTL - segment existing time (in hours)
	// required: false
	// min: 1
	// example: 72
	TTL *int `json:"ttl,omitempty"`
}

type UserSegmentsOutputDTO struct {
	SegmentNames []string `json:"segment_names,omitempty"`
}
