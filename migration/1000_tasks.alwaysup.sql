---
-- This file provides the customized functionality for the session locking package.
---

DROP FUNCTION IF EXISTS get_work(in_session_id session.id%TYPE);
DROP FUNCTION IF EXISTS get_tasks_for_session(in_session_id user_entry.session_id%TYPE);
DROP TYPE IF EXISTS session_task;
CREATE TYPE session_task AS (
    -- TODO - FILL in the info here that you'll need access to in order to "do" the task

    -- user_id     UUID,
    -- stuff       TEXT,
    -- ...

);


-- This will count how many total tasks there are currently to do.
CREATE OR REPLACE FUNCTION get_task_count()
RETURNS INTEGER
AS $$
DECLARE
    v_ret INTEGER;
BEGIN
    -- TODO - Fill in this function so that it returns the total count for current tasks

    -- SELECT count(*) 
    -- FROM task 
    -- INTO v_ret;

    RETURN v_ret;
END;
$$ LANGUAGE plpgsql;

-- This will count how many tasks this session is currently dealing with.
CREATE OR REPLACE FUNCTION get_task_count_for_session(in_session_id user_entry.session_id%TYPE)
RETURNS INTEGER
AS $$
DECLARE
    v_ret INTEGER;
BEGIN
    -- TODO - Fill in this function so that it returns the count for current tasks for the session passed in 
    
    -- SELECT count(*) 
    -- FROM task 
    -- WHERE session_id = in_session_id
    -- INTO v_ret;

    RETURN v_ret;
END;
$$ LANGUAGE plpgsql;

-- This will count how many tasks this session is currently dealing with.
CREATE OR REPLACE FUNCTION pickup_tasks_for_session(in_session_id user_entry.session_id%TYPE
                                                    , in_ideal_pickup INTEGER)
RETURNS VOID
AS $$
DECLARE
    v_now   TIMESTAMP = now() at TIME ZONE 'utc';
BEGIN
    -- TODO - Fill in this function so that it updates N tasks with the session id passed in where N = in_ideal_pickup

    -- UPDATE task t
    -- SET session_id = in_session_id
    -- WHERE t.id = ANY(
    --     SELECT tt.id 
    --     FROM task tt
    --     LEFT OUTER JOIN session s on tt.session_id = s.id
    --     WHERE (tt.session_id IS NULL OR s.expires < v_now) -- there is no session or there is an expired session working this task
    --     LIMIT in_ideal_pickup
    -- );
END;
$$ LANGUAGE plpgsql;

-- This will fetch tasks for a session
CREATE OR REPLACE FUNCTION get_tasks_for_session(in_session_id user_entry.session_id%TYPE)
RETURNS SETOF session_task
AS $$
BEGIN
    -- TODO - Fill in this function so that it returns all tasks this session needs to do

    RETURN QUERY(
        -- SELECT user_id, stuff
        -- FROM task 
        -- WHERE session_id = in_session_id
    );
END;
$$ LANGUAGE plpgsql;

-- This will fetch tasks for a session
CREATE OR REPLACE FUNCTION finish_tasks(in_task_ids BIGINT[])
RETURNS VOID
AS $$
BEGIN
    -- TODO - Fill in this function so that it flags all provided task ids as finished.

    -- UPDATE task SET status = 'finished' WHERE id = ANY(in_task_ids);
END;
$$ LANGUAGE plpgsql;

