import dotenv from 'dotenv';

dotenv.config()
export const BOT_TOKEN = process.env.BOT_TOKEN;
export const BOT_USERNAME = process.env.BOT_USERNAME;

export const time_units = {
    second: 1000,
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000,
};

export const MATCH_NOT_FOUND_TEXT = "Sorry, I don't know what you want. \
Please specify a length of time to remind you after (e.g. 5 Minutes, 2 Days, etc.).";

export const MESSAGE_INCOMPLETE_TEXT = "Time not followed by a valid message. Exiting...";

export const REMINDERS_CLEARED_TEXT = "Reminders cleared.";

// How long to wait for the actual message, if the user sends a time but no message
// This should only happen when adding text to a forward, so should be pretty quick after.
export const MESSAGE_WAIT_TIMEOUT = 10 * time_units.second;

// To preserve real ip address after reverse proxy
export const REAL_IP_HEADER = 'x-real-ip';
