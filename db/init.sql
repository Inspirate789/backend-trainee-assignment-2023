create table if not exists public.users(
    user_id bigint generated always as identity primary key
);

create table if not exists public.segments(
    segment_id bigint generated always as identity primary key,
    name text unique not null
);

create table if not exists public.users_and_segments(
    user_id bigint,
    foreign key (user_id) references public.users(user_id),
    segment_id bigint,
    foreign key (segment_id) references public.segments(segment_id),
    expiration_date timestamp default null
);

-- TODO: triggers on insert and delete for history

create or replace procedure public.insert_user_segments(names text[], user_id bigint, expiration_date timestamp)
as $$
declare segments_id bigint[];
begin
    -- TODO: select into segments_id
    insert into public.users_and_segments(user_id, segment_id, expiration_date)
    values (user_id, unnest(segments_id), expiration_date)
    on conflict do update expiration_date = expiration_date; -- TODO: it works?
end
$$ language sql;

create or replace procedure public.delete_user_segments(names text[], user_id bigint)
as $$
declare segments_id bigint[];
begin
    -- TODO: select into segments_id
    delete from public.users_and_segments us
    where us.user_id = user_id and us.segment_id = any(segments_id);
end
$$ language sql;

create or replace procedure public.clear_expired_linkages()
as $$
begin
    delete from public.users_and_segments
    where expiration_date is not null and expiration_date < now();
end
$$ language sql;

create or replace function public.select_user_segments(user_id bigint)
RETURNS table (
    name text
)
as $$
begin
    public.clear_expired_linkages();
    return query
    select s.name
    from public.segments s join public.users_and_segments us
        on s.segment_id = us.segment_id and us.user_id = user_id;
end
$$ language sql;

create or replace procedure public.insert_segment(name text, user_percentage int, expiration_date timestamp)
as $$
declare segment_id int;
begin
    insert into public.segments(name)
    values (name)
    returning segment_id into segment_id;
    insert into public.users_and_segments(user_id, segment_id, expiration_date)
    values (unnest(select_rand_user_id(user_percentage)), segment_id, expiration_date)
    on conflict do update expiration_date = expiration_date; -- TODO: it works?
    -- TODO: select_rand_user_id
    -- TODO: what if user_percentage is 0?
end
$$ language sql;

create or replace procedure public.delete_segment(name text)
as $$
begin
    delete from public.segments s
    where s.name = name;
end
$$ language sql;

create or replace function delete_users_segment_trigger() -- TODO
returns trigger
as $$
begin
    --RAISE NOTICE 'Old =  %', old;
   	--RAISE NOTICE 'New =  %', new;
--     UPDATE public.users
--     SET last_addition_date = new.last_addition_date;
--     RETURN new; -- Для операций INSERT и UPDATE возвращаемым значением должно быть NEW.
    delete from public.users_and_segments us
    where us.segment_id = old.segment_id;
end;
$$ language sql;

create or replace trigger delete_segment
after delete on public.segments
for each row
execute procedure delete_users_segment_trigger();
