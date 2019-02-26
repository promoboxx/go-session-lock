---
-- This file provides the standard schema for the session locking package
---

-- session holds information about services connected to the db.
-- Each service will start a session when they connect and update the expires timestamp
-- in a keep alive loop.  Once a session expires its work will be given to another session.
CREATE TABLE session (
    id                  BIGSERIAL NOT NULL, 
    created             TIMESTAMP NOT NULL,
    expires             TIMESTAMP NOT NULL,

    CONSTRAINT session_pk1 PRIMARY KEY(id)
);

-- This table is just so we have something to lock during the get_work proc
CREATE TABLE work_lock (
    id                  BIGSERIAL NOT NULL, 
    created             TIMESTAMP NOT NULL,

    CONSTRAINT work_lock_pk1 PRIMARY KEY(id)
);
INSERT INTO work_lock (id, created) VALUES(1, now() at TIME ZONE 'utc');


-- TODO - EDIT below this line to add a session_id BIGINT column to the table that is keeping track of tasks to do

