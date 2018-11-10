'use strict';

var _http = require('http');

var _http2 = _interopRequireDefault(_http);

var _express = require('express');

var _express2 = _interopRequireDefault(_express);

var _bodyParser = require('body-parser');

var _bodyParser2 = _interopRequireDefault(_bodyParser);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _resources = require('./resources');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// Users that have forwarded messages to the bot, but havent sent a time yet, or visa versa
var USERS_WAITING = {};
// ------------------------------------------------------

// import https from 'https';
var app = (0, _express2.default)();
app.use(_bodyParser2.default.json());
app.use(_bodyParser2.default.urlencoded({
    extended: true
}));

app.get('/', function (req, res) {
    console.log(new Date().toLocaleString() + '.....' + 'GET REQUEST from: ' + req.headers[_resources.REAL_IP_HEADER]);

    res.status(403).send('Knock it off');
    res.end();
});

// This is the route the API will call
app.post('/', function (req, res) {
    res.end();

    if (!(0, _resources.checkIfFromTelegram)(req)) {
        console.log('NOT FROM TELEGRAM, quitting...');
        return;
    }

    var message = req.body.message;
    if (!message) {
        console.log("No message!");
        return;
    }

    // Message information --------------
    var sender_id = message.from.id;
    var message_id = message.message_id;
    var chat_id = message.chat.id;

    var match = (0, _resources.findMatchFromText)(message.text);

    var REPLY = message.reply_to_message ? true : false;
    var FORWARD = message.forward_from_message_id ? true : false;
    var PRIVATE = message.chat.type === 'private' ? true : false;
    // ----------------------------------

    var params = {
        chat_id: sender_id,
        from_chat_id: chat_id
    };

    if (REPLY) params.message_id = message.reply_to_message.message_id;else if (FORWARD) params.message_id = message.forward_from_message_id;else params.message_id = message_id;

    var SHOULD_WAIT = match.found && !REPLY && !FORWARD;
    if (PRIVATE) {
        // A message was sent directly to the bot.
        console.log('private chat');

        if (USERS_WAITING[sender_id] !== undefined) {
            // Message from someone we're waiting for

            match = USERS_WAITING[sender_id];
            delete USERS_WAITING[sender_id];

            params.assembled = true;
        } else if (SHOULD_WAIT) {
            // We've only recieved half of what we need, so wait.

            USERS_WAITING[sender_id] = match;
            setTimeout(function (params) {
                if (USERS_WAITING[params.sender_id]) {
                    (0, _resources.sendMessage)({
                        chat_id: params.chat_id,
                        reply_to_message_id: params.message_id,
                        text: _resources.MESSAGE_INCOMPLETE_TEXT
                    });
                }
                delete USERS_WAITING[params.sender_id];
            }, _resources.MESSAGE_WAIT_TIMEOUT, { sender_id: sender_id, message_id: message_id, chat_id: chat_id });
        }
    } else {
        // Public, bot must be mentioned
        console.log('public chat');

        if (!_lodash2.default.some(_resources.KEYWORDS, function (x) {
            return _lodash2.default.toLower(message.text).includes(_lodash2.default.toLower(x));
        })) {
            // Bot not mentioned, ignore

            console.log("Bot not mentioned, let's not get this bread.");
            return;
        }
    }

    // We've gotten here, which means that the bot either had a message sent directly to it, or it has been mentioned publicly
    if (!SHOULD_WAIT) {
        if (match.found) {
            setTimeout(function () {
                (0, _resources.forwardMessage)(params);
            }, match.wait);

            // Confirm to sender that reminder is set.
            var text = 'Reminder set for ' + match.num + ' ' + match.units + ' from now.';
            if (!REPLY && !FORWARD && !params.assembled) text = 'No message specified. ' + text;

            (0, _resources.sendMessage)({
                chat_id: chat_id,
                reply_to_message_id: message_id,
                text: text
            });
        } else {
            (0, _resources.sendMessage)({
                chat_id: chat_id,
                reply_to_message_id: message_id,
                text: _resources.MATCH_NOT_FOUND_TEXT
            });
        }
    }
});

_http2.default.createServer(app).listen(8081);
console.log('Listening on port 443.....');