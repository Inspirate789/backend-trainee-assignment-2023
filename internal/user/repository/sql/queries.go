package sql

const (
	insertUserQuery         = `insert into public.users(id) values (:user_id) on conflict do nothing;`
	deleteUserQuery         = `delete from public.users u where u.id = :user_id;`
	selectUserSegmentsQuery = `select public.select_user_segments(:user_id);`
	insertUserSegmentsQuery = `call public.insert_user_segments(:names, :user_id, :expire);`
	deleteUserSegmentsQuery = `call public.delete_user_segments(:names, :user_id);`
	selectUserHistoryQuery  = `select * from public.user_history uh 
         where extract(year from uh.registration_date) = :year and extract(month from uh.registration_date) = :month;`
)
