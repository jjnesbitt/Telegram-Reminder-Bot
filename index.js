import https from 'https';
import http from 'http';
import express from 'express';
import bodyParser from 'body-parser';
import axios from 'axios';
import fs from 'fs';

import { BOT_TOKEN, BOT_USERNAME } from './credentials';

let app = express();

process.exit(0);

const MATCH_NOT_FOUND_TEXT = "Sorry, I don't know what you want."
const SSL_OPTIONS = {
    key: fs.readFileSync('./telegram_private.key'),
    cert: fs.readFileSync('./telegram_public.pem')
};

//In milliseconds
const time_units = {
    second: 1000,
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000
};

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({
    extended: true
})); // for parsing application/x-www-form-urlencoded

//Function called to forward message, usually set on a timer.
function forwardMessage(params) {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/forwardMessage', params)
        .then(response => {
            // We get here if the message was successfully posted
            console.log('Message forwarded');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err);
        });
}

//Function called to send message, usually set on a timer.
function sendMessage(params) {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/sendMessage', params)
        .then(response => {
            // We get here if the message was successfully posted
            console.log('Message sent');
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err);
        });
}

function findMatchFromText(text) {
    const found = 'found';

    let returnObject = {};

    if (!text) {
        returnObject[found] = false;
    }
    else {
        let matches = [];
        for (let key in time_units) {
            //let regex = new RegExp(BOT_USERNAME + " (\\d{1,10}.?\\d{0,10}) " + key);
            let regex = new RegExp("(\\d{1,10}.?\\d{0,10}) " + key);
            let match;
            if ((match = regex.exec(text)) != null) {
                //Means valid number has been matched

                let num = match[1];
                matches.push(parseFloat(num));
                matches.push(key);
                break;
            }
        }

        if (matches.length == 0) {
            returnObject[found] = false;
        }
        else {
            let num = matches[0];
            let units = matches[matches.length - 1];
            let timeToWait = num * time_units[units];

            console.log("TTW: " + timeToWait);

            //Format before return
            if (num != 1) units += 's';

            returnObject[found] = true;
            returnObject['wait'] = timeToWait;
            returnObject['num'] = num;
            returnObject['units'] = units;
        }
    }

    return returnObject;
}

function confirmReminderSet(chat_id, reply_to_message_id, numberOfUnits, units) {
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/sendMessage', {
        'chat_id': chat_id,
        'text': ('Reminder set for ' + numberOfUnits + ' ' + units + ' from now.'),
        'reply_to_message_id': reply_to_message_id
    })
        .then(response => {
            // We get here if the message was successfully posted
            console.log('Message posted')
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err)
        })
}

function checkIfFromTelegram(req){
    let ip = req.connection.remoteAddress.split('::ffff:')[1];

    console.log('REQUEST IP: ' + ip);
    const ipToCheck = '149.154.167';
    let splitIP = ip.split('.');

    if (splitIP[0] + '.' + splitIP[1] + '.' + splitIP[2]  != ipToCheck) return false;
    if (parseInt(splitIP[3]) < 197 || parseInt(splitIP[3]) > 233) return false;

    return true;
}

app.get('/', function (req, res) {
    console.log(new Date().toLocaleString() + '.....' + 'GET REQUEST from: ' + req.connection.remoteAddress.split('::ffff:')[1]);

    res.status(403).send('Knock it off');
    res.end();
});

//This is the route the API will call
app.post('/', function (req, res) {
    res.end();

    if (!checkIfFromTelegram(req)){
        console.log('NOT FROM TELEGRAM, quitting...');
        return;
    }

    const message = req.body.message;

    if (!message) {
        console.log("No message!");
        return;
    }

    let sender_id = message.from.id;
    let message_id = message.message_id;
    let chat_id = message.chat.id;

    let match = findMatchFromText(message.text);

    let messageFunction;
    let params = {
        chat_id: sender_id,
        from_chat_id: chat_id
    };

    let REPLY = false;
    let FORWARD = false;

    messageFunction = forwardMessage;
    // params['from_chat_id'] = chat_id;

    if (message.chat.type == 'private') {
        //A message was sent directly to the bot.
        //Will forward message back to user

        console.log('private chat');

        if (message.forward_from_message_id){
            //Contains a forwarded message
            FORWARD = true;
            params['message_id'] = message.forward_from_message_id;
        }
        else if (message.reply_to_message){
            //Is a reply to a message
            REPLY = true;
            params['message_id'] = message.reply_to_message.message_id;
        }
        else{
            params['message_id'] = message_id;
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
                params['message_id'] = message.reply_to_message.message_id;
            }
            else {
                params['message_id'] = message_id;
            }
        }
        else{
            //Bot not mentioned, ignore
            return;
        }
    }

    //We've gotten here, which means that the bot either had a message sent directly to it, or it has been mentioned publicly
    if (match.found){
        //Either forward or send message
        setTimeout(function(){messageFunction(params);}, match.wait);

        //Confirm to sender that reminder is set.
        let text = 'Reminder set for ' + match.num + ' ' + match.units + ' from now.';
        if (!REPLY && !FORWARD) text = "No message specified, " + text;

        sendMessage({
            'chat_id': chat_id,
            'reply_to_message_id': message_id,
            'text': text
        });
    }
    else{
        sendMessage({
            'chat_id': chat_id,
            'reply_to_message_id': message_id,
            'text': MATCH_NOT_FOUND_TEXT
        });
    }
    return;
});

//create server to listen for api call
https.createServer(SSL_OPTIONS, app).listen(443);
console.log('Listening on port 443.....');

http.createServer(app).listen(80);
console.log('Listening on port 80.....');
