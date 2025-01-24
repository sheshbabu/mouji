ALTER TABLE pageviews
    ADD COLUMN visitor_hash TEXT;

ALTER TABLE pageviews
    ADD COLUMN user_agent TEXT;
