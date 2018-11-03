// import https from 'https';
import http from 'http';
import express from 'express';
import bodyParser from 'body-parser';
import axios from 'axios';
// import fs from 'fs';
import _ from 'lodash';

import { BOT_TOKEN, BOT_USERNAME } from './credentials';

const REAL_IP_HEADER = 'x-real-ip';

const app = express();

const MATCH_NOT_FOUND_TEXT = "Sorry, I don't know what you want. \
Please specify a length of time to remind you after (e.g. 5 Minutes, 2 Days, etc.).";
// const SSL_OPTIONS = {
//     key: fs.readFileSync('./telegram_private.key'),
//     cert: fs.readFileSync('./telegram_public.pem'),
// };

//In milliseconds
const time_units = {
    second: 1000,
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000,
};


app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({
    extended: true,
})); // for parsing application/x-www-form-urlencoded

//Function called to forward message, usually set on a timer.

const forwardMessage = params => {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/forwardMessage', params)
        .then(() => {
            // We get here if the message was successfully posted
            console.log('Message forwarded');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err);
        });
};

//Function called to send message, usually set on a timer.
function sendMessage(params) {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/sendMessage', params)
        .then(() => {
            // We get here if the message was successfully posted
            console.log('Message sent');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err);
        });
}

function findMatchFromText(text) {
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
}

function checkIfFromTelegram(req) {
    const ipToCheck = '149.154.167';
    
    // const ipArray = req.connection.remoteAddress.split('::ffff:')[1].split('.');
    const ipArray = req.headers[REAL_IP_HEADER].split('.');
    console.log('REQUEST IP: ' + ipArray.join('.'));

    // if (ipArray[0] + '.' + ipArray[1] + '.' + ipArray[2] != ipToCheck) return false;
    if (_.dropRight(ipArray).join('.') != ipToCheck) return false;
    if (parseInt(ipArray[3]) < 197 || parseInt(ipArray[3]) > 233) return false;

    return true;
}

app.get('/', function (req, res) {
    console.log(new Date().toLocaleString() + '.....' + 'GET REQUEST from: ' + req.headers[REAL_IP_HEADER]);

    res.status(403).send('Knock it off');
    res.end();
});

//This is the route the API will call
app.post('/', (req, res) => {
    res.end();

    if (!checkIfFromTelegram(req)) {
        console.log('NOT FROM TELEGRAM, quitting...');
        return;
    }

    const message = req.body.message;

    if (!message) {
        console.log("No message!");
        return;
    }

    const sender_id = message.from.id;
    const message_id = message.message_id;
    const chat_id = message.chat.id;

    const match = findMatchFromText(message.text);

    const params = {
        chat_id: sender_id,
        from_chat_id: chat_id,
    };

    let REPLY = false;
    let FORWARD = false;

    // params['from_chat_id'] = chat_id;

    if (message.chat.type == 'private') {
        //A message was sent directly to the bot.
        //Will forward message back to user

        console.log('private chat');

        if (message.forward_from_message_id) {
            //Contains a forwarded message
            FORWARD = true;
            params.message_id = message.forward_from_message_id;
        }
        else if (message.reply_to_message) {
            //Is a reply to a message
            REPLY = true;
            params.message_id = message.reply_to_message.message_id;
        }
        else {
            params.message_id = message_id;
        }
    }
    else {
        //Public, bot must be mentioned
        console.log('public chat');
        if (RegExp(BOT_USERNAME).test(message.text)) {
            //Bot mentioned in a group, continue...

            if (message.reply_to_message) {
                //Message is a reply to another message
                REPLY = true;
                params.message_id = message.reply_to_message.message_id;
            }
            else {
                params.message_id = message_id;
            }
        }
        else {
            //Bot not mentioned, ignore
            console.log("Bot not mentioned, let's not get this bread.");
            return;
        }
    }

    //We've gotten here, which means that the bot either had a message sent directly to it, or it has been mentioned publicly
    if (match.found) {
        //Either forward or send message
        setTimeout(() => { forwardMessage(params); }, match.wait);

        //Confirm to sender that reminder is set.
        let text = 'Reminder set for ' + match.num + ' ' + match.units + ' from now.';
        if (!REPLY && !FORWARD) text = "No message specified, " + text;

        sendMessage({
            chat_id,
            text,
            reply_to_message_id: message_id,
        });
    }
    else {
        sendMessage({
            chat_id,
            reply_to_message_id: message_id,
            text: MATCH_NOT_FOUND_TEXT,
        });
    }
    return;
});

// Must fix networking with nginx past here

// //create server to listen for api call  
// https.createServer(SSL_OPTIONS, app).listen(8081);
// console.log('Listening on port 443.....');

// https.createServer(app).listen(8081);
// console.log('Listening on port 443.....');

http.createServer(app).listen(8081);
console.log('Listening on port 443.....');
