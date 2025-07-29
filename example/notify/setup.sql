-- Create a test table
CREATE TABLE IF NOT EXISTS notification_test (
    id SERIAL PRIMARY KEY,
    message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create a function that sends notifications
CREATE OR REPLACE FUNCTION notify_new_message() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('test_channel', 
        'New message: ' || NEW.message || ' (ID: ' || NEW.id || ')');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger that fires on INSERT
DROP TRIGGER IF EXISTS new_message_notify ON notification_test;
CREATE TRIGGER new_message_notify
    AFTER INSERT ON notification_test
    FOR EACH ROW
    EXECUTE FUNCTION notify_new_message();