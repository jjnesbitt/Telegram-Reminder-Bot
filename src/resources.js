import axios from 'axios';
import _ from 'lodash';
import { BOT_TOKEN, BOT_USERNAME } from './credentials';

// Exports -------------------

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

// Keywords to invoke the bot
export const KEYWORDS = ['!remindme', BOT_USERNAME];

// How long to wait for the actual message, if the user sends a time but no message
// This should only happen when adding text to a forward, so should be pretty quick after.
export const MESSAGE_WAIT_TIMEOUT = 10*time_units.second;

// To preserve real ip address after reverse proxy
export const REAL_IP_HEADER = 'x-real-ip';

export const forwardMessage = params => {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/forwardMessage', params)
        .then(() => {
            // We get here if the message was successfully posted
            console.log('Message forwarded');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error:', err.keys());
        });
};

// Function called to send message, usually set on a timer.
export const sendMessage = params => {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/sendMessage', params)
        .then(() => {
            // We get here if the message was successfully posted
            console.log('Message sent');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error:', err.keys());
        });
};

export const findMatchFromText = text => {
    // const COMMAND_REGEX = /(\/[^/\s]+\s*)(.+)?/;
    // const matched = COMMAND_REGEX.exec(raw_text);

    const regex = new RegExp(`(\\d{1,10}(?:\\.\\d{0,10})?) (${_.keys(time_units).join('|')})`);
    const match = regex.exec(text);

    if (!text || !match || !match.length) return { found: false };

    const num = parseFloat(match[1]);
    const unit = match[2];
    const wait = num * time_units[unit];

    console.log("TTW: " + wait);

    return {
        wait,
        num,
        found: true,
        units: num === 1 ? unit : `${unit}s`,
    };
};

export const checkIfFromTelegram = req => {
    const REAL_IP_HEADER = 'x-real-ip';
    const ipToCheck = '149.154.167';

    // const ipArray = req.connection.remoteAddress.split('::ffff:')[1].split('.');
    const ipArray = req.headers[REAL_IP_HEADER].split('.');
    // console.log('REQUEST IP: ' + ipArray.join('.'));

    // if (ipArray[0] + '.' + ipArray[1] + '.' + ipArray[2] != ipToCheck) return false;
    if (_.dropRight(ipArray).join('.') != ipToCheck) return false;
    if (parseInt(ipArray[3]) < 197 || parseInt(ipArray[3]) > 233) return false;

    return true;
};