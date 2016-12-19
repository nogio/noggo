CREATE OR REPLACE FUNCTION jsonb_num(data jsonb, name varchar, number int8)
	RETURNS jsonb AS
$BODY$
    DECLARE
		value jsonb;
		count int8;
    BEGIN
		SELECT data->name into count;
		IF count is null THEN
			SELECT 0 into count;
		END IF;
		SELECT count+number into count;

		select jsonb_set(data,('{'||name||'}')::text[], (count::varchar)::jsonb, true) into value;


        RETURN value;
    END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
