---
-- This file provides the standard functionality for the session locking package
---

---
CREATE OR REPLACE FUNCTION throw_session_not_found()
RETURNS VOID 
AS $$
BEGIN
    RAISE EXCEPTION 'Session not found.' USING ERRCODE = 'SL001';
END;
$$ LANGUAGE plpgsql;

---
-- This will start a new session for a service.
---
CREATE OR REPLACE FUNCTION start_session()
RETURNS BIGINT
AS $$
DECLARE
    v_ret BIGINT;
    v_now TIMESTAMP = now() at TIME ZONE 'utc';
BEGIN
    INSERT INTO session (created, expires) VALUES (v_now, v_now + INTERVAL '2 minutes') RETURNING id INTO v_ret;
    RETURN v_ret;
END;
$$ LANGUAGE plpgsql;

---
-- This will update an existing session for a service to keep it active
---
CREATE OR REPLACE FUNCTION bump_session(in_session_id session.id%TYPE)
RETURNS VOID
AS $$
DECLARE
    v_now TIMESTAMP = now() at TIME ZONE 'utc';
BEGIN
    UPDATE session
    SET expires = v_now + INTERVAL '2 minutes'
    WHERE id = in_session_id
    AND expires >= v_now;

    IF NOT FOUND THEN
        PERFORM throw_session_not_found();
    END IF;
END;
$$ LANGUAGE plpgsql;


---
-- This will end a session and return it's tasks to the pull
---
CREATE OR REPLACE FUNCTION end_session(in_session_id session.id%TYPE)
RETURNS VOID
AS $$
DECLARE
    v_now TIMESTAMP = now() at TIME ZONE 'utc';
BEGIN
    UPDATE session
    SET expires = '-INFINITY'
    WHERE id = in_session_id;
END;
$$ LANGUAGE plpgsql;


---
-- This will balance the tasks evenly across the active sessions and 
-- return work for this session to do.
---
CREATE OR REPLACE FUNCTION get_work(in_session_id session.id%TYPE)
RETURNS SETOF session_task
AS $$
DECLARE
    v_now           TIMESTAMP = now() at TIME ZONE 'utc';
    v_sessions      INTEGER;
    v_task_count    INTEGER;
    v_ideal_count   INTEGER;
    v_session_count INTEGER;
BEGIN
    -- lock work
    PERFORM 1 FROM work_lock WHERE id = 1 FOR UPDATE;

    -- bump this session to extend its expiration time
    PERFORM bump_session(in_session_id);
    
    -- count active sessions
    SELECT count(*) FROM session WHERE expires >= v_now INTO v_sessions;
    -- count active tasks and calculate ideal task count per session (rounded up)
    SELECT get_task_count FROM get_task_count() INTO v_task_count;
    v_ideal_count := CEIL(v_task_count::NUMERIC / v_sessions::NUMERIC)::INTEGER;
    -- count how many active tasks this session has
    SELECT get_task_count_for_session FROM get_task_count_for_session(in_session_id) INTO v_session_count;

    -- distribute tasks - i.e. pickup unassociated tasks if necessary
    IF v_session_count < v_ideal_count THEN 
        -- pick up tasks if possible     
        PERFORM pickup_tasks_for_session(in_session_id, v_ideal_count - v_session_count);
    END IF;

    -- return tasks that are ready to run
    RETURN QUERY (
        SELECT * FROM get_tasks_for_session(in_session_id)
    );
END;
$$ LANGUAGE plpgsql;



