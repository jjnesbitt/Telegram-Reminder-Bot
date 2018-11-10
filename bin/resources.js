'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.checkIfFromTelegram = exports.checkCommands = exports.findMatchFromText = exports.sendMessage = exports.forwardMessage = exports.updateReminders = exports.cancelReminders = exports.REAL_IP_HEADER = exports.MESSAGE_WAIT_TIMEOUT = exports.KEYWORDS = exports.REMINDERS_CLEARED_TEXT = exports.MESSAGE_INCOMPLETE_TEXT = exports.MATCH_NOT_FOUND_TEXT = exports.time_units = undefined;

var _axios = require('axios');

var _axios2 = _interopRequireDefault(_axios);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _credentials = require('./credentials');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

// Exports -------------------

var time_units = exports.time_units = {
    second: 1000,
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000
};

// Object that stores the timout ID's for forwarding messages, by user id
var currentReminders = {};

var MATCH_NOT_FOUND_TEXT = exports.MATCH_NOT_FOUND_TEXT = "Sorry, I don't know what you want. \
Please specify a length of time to remind you after (e.g. 5 Minutes, 2 Days, etc.).";

var MESSAGE_INCOMPLETE_TEXT = exports.MESSAGE_INCOMPLETE_TEXT = "Time not followed by a valid message. Exiting...";

var REMINDERS_CLEARED_TEXT = exports.REMINDERS_CLEARED_TEXT = "Reminders cleared.";

// Keywords to invoke the bot
var KEYWORDS = exports.KEYWORDS = ['!remindme', _credentials.BOT_USERNAME];

// How long to wait for the actual message, if the user sends a time but no message
// This should only happen when adding text to a forward, so should be pretty quick after.
var MESSAGE_WAIT_TIMEOUT = exports.MESSAGE_WAIT_TIMEOUT = 10 * time_units.second;

// To preserve real ip address after reverse proxy
var REAL_IP_HEADER = exports.REAL_IP_HEADER = 'x-real-ip';

var cancelReminders = exports.cancelReminders = function cancelReminders(sender_id) {
    var timeouts = currentReminders[sender_id];
    if (timeouts !== undefined) {
        timeouts.forEach(function (val) {
            return clearTimeout(val);
        });
        delete currentReminders[sender_id];
    }
};

var updateReminders = exports.updateReminders = function updateReminders(sender_id, timeout_id) {
    currentReminders[sender_id] = [].concat(_toConsumableArray(currentReminders[sender_id]), [timeout_id]);
};

var forwardMessage = exports.forwardMessage = function forwardMessage(params) {
    _axios2.default.post('https://api.telegram.org/bot' + _credentials.BOT_TOKEN + '/forwardMessage', params).then(function () {
        // We get here if the message was successfully posted
        console.log('Message forwarded');
    }).catch(function (err) {
        // ...and here if it was not
        console.log('Error:', err.keys());
    });
};

// Function called to send message, usually set on a timer.
var sendMessage = exports.sendMessage = function sendMessage(params) {
    _axios2.default.post('https://api.telegram.org/bot' + _credentials.BOT_TOKEN + '/sendMessage', params).then(function () {
        // We get here if the message was successfully posted
        console.log('Message sent');
    }).catch(function (err) {
        // ...and here if it was not
        console.log('Error:', err.keys());
    });
};

var findMatchFromText = exports.findMatchFromText = function findMatchFromText(text) {
    // const COMMAND_REGEX = /(\/[^/\s]+\s*)(.+)?/;
    // const matched = COMMAND_REGEX.exec(raw_text);

    var regex = new RegExp('(\\d{1,10}(?:\\.\\d{0,10})?) (' + _lodash2.default.keys(time_units).join('|') + ')');
    var match = regex.exec(text);

    if (!text || !match || !match.length) return { found: false };

    var num = parseFloat(match[1]);
    var unit = match[2];
    var wait = num * time_units[unit];

    console.log("TTW: " + wait);

    return {
        wait: wait,
        num: num,
        found: true,
        units: num === 1 ? unit : unit + 's'
    };
};

var checkCommands = exports.checkCommands = function checkCommands(message) {
    var matchCommandRegex = /^\/(\w+)/;
    var match = matchCommandRegex.exec(message.text);

    console.log(match);
    if (!match) return false;

    var command = match[1];
    if (command == 'start') return true;else if (command == 'cancel') {
        cancelReminders(message.from.id);
        sendMessage({
            chat_id: message.chat.id,
            reply_to_message_id: message.message_id,
            text: REMINDERS_CLEARED_TEXT
        });
        return true;
    } else if (command == 'reminders') {
        console.log("current reminders", currentReminders[message.from.id]);

        // finish later

        // sendMessage({
        //     chat_id: message.chat.id,
        //     reply_to_message_id: message.message_id,
        //     text: REMINDERS_CLEARED_TEXT,
        // });

        // return false for now cuz this command isn't implemented yet.
        // return true;
        return false;
    }

    return false;
};

var checkIfFromTelegram = exports.checkIfFromTelegram = function checkIfFromTelegram(req) {
    var REAL_IP_HEADER = 'x-real-ip';
    var ipToCheck = '149.154.167';

    // const ipArray = req.connection.remoteAddress.split('::ffff:')[1].split('.');
    var ipArray = req.headers[REAL_IP_HEADER].split('.');
    // console.log('REQUEST IP: ' + ipArray.join('.'));

    // if (ipArray[0] + '.' + ipArray[1] + '.' + ipArray[2] != ipToCheck) return false;
    if (_lodash2.default.dropRight(ipArray).join('.') != ipToCheck) return false;
    if (parseInt(ipArray[3]) < 197 || parseInt(ipArray[3]) > 233) return false;

    return true;
};