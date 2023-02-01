CREATE OR REPLACE FUNCTION fun_setrole() RETURNS TRIGGER AS
$BODY$
BEGIN
    insert into l_role_user (user_id, role_id) values (new.id, (select id from roles where name='user'));
    return new;
END;
$BODY$
    language plpgsql;

CREATE TRIGGER TRI_USERS
    AFTER INSERT
    ON USERS
    FOR EACH ROW
EXECUTE PROCEDURE fun_setrole();