package repository

const (
	insertSegmentQuery = `call public.insert_segment(:seg_name, :user_percentage, :expire);`
	deleteSegmentQuery = `call public.delete_segment(:seg_name);`
)
