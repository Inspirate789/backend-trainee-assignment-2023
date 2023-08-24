create table if not exists public.segments(
    segment_id int generated always as identity primary key,
    name text unique not null
);

create table if not exists public.users_and_segments(
    user_id int,
    segment_id int,
    foreign key (segment_id) references public.segments(segment_id),
    registration_date timestamp not null default now(),
    expiration_date timestamp default null
);

create or replace procedure public.clear_expired_linkages()
AS $$
BEGIN
    delete from public.users_and_segments
    where expiration_date is not null and expiration_date < now();
END
$$ LANGUAGE sql;

create or replace function public.select_user_segments(user_id int)
RETURNS table (
    name text
)
AS $$
BEGIN
    public.clear_expired_linkages();

    return query
    select s.name
    from public.segments s join public.users_and_segments us
        on s.segment_id = us.segment_id and us.user_id = user_id;
END
$$ LANGUAGE sql;

create or replace procedure public.insert_segment(name text, user_percentage int)
AS $$
DECLARE segment_id int;
BEGIN
    insert into public.segments(name)
    values (name)
    returning segment_id into segment_id;

    insert into public.users_and_segments(user_id, segment_id)
    values (unnest(select_rand_user_id(user_percentage)), segment_id); -- TODO: select_rand_user_id
    -- TODO: expiration date?
END
$$ LANGUAGE sql;
