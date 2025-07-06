CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION public.generate_uuid_v7()
RETURNS uuid
LANGUAGE plpgsql
AS $function$
DECLARE
    ts_millis BIGINT;
    rand_bytes BYTEA;
    uuid_hex TEXT;
BEGIN
    -- Get current timestamp in milliseconds (48 bits)
    ts_millis := (extract(epoch FROM now()) * 1000)::BIGINT;

    -- Generate 10 random bytes (80 bits)
    rand_bytes := gen_random_bytes(10);

    -- Construct the UUID manually
    uuid_hex := lpad(to_hex(ts_millis), 12, '0') ||
                substr(encode(rand_bytes, 'hex'), 1, 4) ||
                '7' || substr(encode(rand_bytes, 'hex'), 6, 3) ||
                substr(encode(rand_bytes, 'hex'), 9, 4) ||
                substr(encode(rand_bytes, 'hex'), 13, 12);

    RETURN uuid_hex::UUID;
END;
$function$;
