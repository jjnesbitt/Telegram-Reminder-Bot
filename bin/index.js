'use strict';

var _http = require('http');

var _http2 = _interopRequireDefault(_http);

var _express = require('express');

var _express2 = _interopRequireDefault(_express);

var _bodyParser = require('body-parser');

var _bodyParser2 = _interopRequireDefault(_bodyParser);

var _axios = require('axios');

var _axios2 = _interopRequireDefault(_axios);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _credentials = require('./credentials');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// import fs from 'fs';
// import https from 'https';
var REAL_IP_HEADER = 'x-real-ip';
var MATCH_NOT_FOUND_TEXT = "Sorry, I don't know what you want. \
Please specify a length of time to remind you after (e.g. 5 Minutes, 2 Days, etc.).";
var KEYWORDS = ['!remindme'];

var app = (0, _express2.default)();

// In milliseconds
var time_units = {
    second: 1000,
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000
};

app.use(_bodyParser2.default.json()); // for parsing application/json
app.use(_bodyParser2.default.urlencoded({
    extended: true
})); // for parsing application/x-www-form-urlencoded

// Function called to forward message, usually set on a timer.

var forwardMessage = function forwardMessage(params) {
    _axios2.default.post('https://api.telegram.org/bot' + _credentials.BOT_TOKEN + '/forwardMessage', params).then(function () {
        // We get here if the message was successfully posted
        console.log('Message forwarded');
    }).catch(function (err) {
        // ...and here if it was not
        console.log('Error :', err);
    });
};

//Function called to send message, usually set on a timer.
function sendMessage(params) {
    _axios2.default.post('https://api.telegram.org/bot' + _credentials.BOT_TOKEN + '/sendMessage', params).then(function () {
        // We get here if the message was successfully posted
        console.log('Message sent');
    }).catch(function (err) {
        // ...and here if it was not
        console.log('Error :', err);
    });
}

function findMatchFromText(text) {
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
}

function checkIfFromTelegram(req) {
    var ipToCheck = '149.154.167';

    // const ipArray = req.connection.remoteAddress.split('::ffff:')[1].split('.');
    var ipArray = req.headers[REAL_IP_HEADER].split('.');
    // console.log('REQUEST IP: ' + ipArray.join('.'));

    // if (ipArray[0] + '.' + ipArray[1] + '.' + ipArray[2] != ipToCheck) return false;
    if (_lodash2.default.dropRight(ipArray).join('.') != ipToCheck) return false;
    if (parseInt(ipArray[3]) < 197 || parseInt(ipArray[3]) > 233) return false;

    return true;
}

app.get('/', function (req, res) {
    console.log(new Date().toLocaleString() + '.....' + 'GET REQUEST from: ' + req.headers[REAL_IP_HEADER]);

    res.status(403).send('Knock it off');
    res.end();
});

//This is the route the API will call
app.post('/', function (req, res) {
    res.end();

    if (!checkIfFromTelegram(req)) {
        console.log('NOT FROM TELEGRAM, quitting...');
        return;
    }

    var message = req.body.message;

    if (!message) {
        console.log("No message!");
        return;
    }

    var sender_id = message.from.id;
    var message_id = message.message_id;
    var chat_id = message.chat.id;

    var match = findMatchFromText(message.text);

    var params = {
        chat_id: sender_id,
        from_chat_id: chat_id
    };

    var REPLY = false;
    var FORWARD = false;

    // params['from_chat_id'] = chat_id;

    if (message.chat.type == 'private') {
        //A message was sent directly to the bot.
        //Will forward message back to user

        console.log('private chat');

        if (message.forward_from_message_id) {
            //Contains a forwarded message
            FORWARD = true;
            params.message_id = message.forward_from_message_id;
        } else if (message.reply_to_message) {
            //Is a reply to a message
            REPLY = true;
            params.message_id = message.reply_to_message.message_id;
        } else {
            params.message_id = message_id;
        }
    } else {
        //Public, bot must be mentioned
        console.log('public chat');

        // if (RegExp(BOT_USERNAME).test(message.text)) {
        if (_lodash2.default.some([_credentials.BOT_USERNAME].concat(KEYWORDS), function (x) {
            return _lodash2.default.toLower(message.text).includes(_lodash2.default.toLower(x));
        })) {
            //Bot mentioned in a group, continue...

            if (message.reply_to_message) {
                //Message is a reply to another message
                REPLY = true;
                params.message_id = message.reply_to_message.message_id;
            } else {
                params.message_id = message_id;
            }
        } else {
            //Bot not mentioned, ignore
            console.log("Bot not mentioned, let's not get this bread.");
            return;
        }
    }

    //We've gotten here, which means that the bot either had a message sent directly to it, or it has been mentioned publicly
    if (match.found) {
        //Either forward or send message
        setTimeout(function () {
            forwardMessage(params);
        }, match.wait);

        //Confirm to sender that reminder is set.
        var text = 'Reminder set for ' + match.num + ' ' + match.units + ' from now.';
        if (!REPLY && !FORWARD) text = 'No message specified. ' + text;

        sendMessage({
            chat_id: chat_id,
            text: text,
            reply_to_message_id: message_id
        });
    } else {
        sendMessage({
            chat_id: chat_id,
            reply_to_message_id: message_id,
            text: MATCH_NOT_FOUND_TEXT
        });
    }
    return;
});

_http2.default.createServer(app).listen(8081);
console.log('Listening on port 443.....');