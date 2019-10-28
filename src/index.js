// import https from 'https';
import http from 'http';
import express from 'express';
import bodyParser from 'body-parser';
import _ from 'lodash';

import {
  forwardMessage,
  sendMessage,
  findMatchFromText,
  checkIfFromTelegram,
  updateReminders,
  checkCommands,
} from './utils';

import { REAL_IP_HEADER, MATCH_NOT_FOUND_TEXT, MESSAGE_INCOMPLETE_TEXT, MESSAGE_WAIT_TIMEOUT } from './constants';

// Users that have forwarded messages to the bot, but havent sent a time yet, or visa versa
const USERS_WAITING = {};

// ------------------------------------------------------

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
  extended: true,
}));

app.get('/', function (req, res) {
  console.log(new Date().toLocaleString() + '.....' + 'GET REQUEST from: ' + req.headers[REAL_IP_HEADER]);

  res.status(405).send('Knock it off');
  res.end();
});

// This is the route the API will call
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

  // Message information --------------
  const sender_id = message.from.id;
  const message_id = message.message_id;
  const chat_id = message.chat.id;

  let match = findMatchFromText(message.text);

  const REPLY = message.reply_to_message ? true : false;
  const FORWARD = message.forward_from_message_id ? true : false;
  const PRIVATE = message.chat.type === 'private' ? true : false;
  // ----------------------------------

  const params = {
    chat_id: sender_id,
    from_chat_id: chat_id,
  };

  // Check if it's a command, and handle
  if (checkCommands(message)) return;

  if (REPLY) params.message_id = message.reply_to_message.message_id;
  else if (FORWARD) params.message_id = message.forward_from_message_id;
  else params.message_id = message_id;

  const SHOULD_WAIT = (match.found && !REPLY && !FORWARD);
  if (PRIVATE) {
    // A message was sent directly to the bot.
    console.log('private chat');

    if (USERS_WAITING[sender_id] !== undefined) {
      // Message from someone we're waiting for

      match = USERS_WAITING[sender_id];
      delete USERS_WAITING[sender_id];

      params.assembled = true;
    }
    else if (SHOULD_WAIT) {
      // We've only recieved half of what we need, so wait.

      USERS_WAITING[sender_id] = match;
      setTimeout((params) => {
        if (USERS_WAITING[params.sender_id]) {
          sendMessage({
            chat_id: params.chat_id,
            reply_to_message_id: params.message_id,
            text: MESSAGE_INCOMPLETE_TEXT,
          });
        }
        delete USERS_WAITING[params.sender_id];
      }, MESSAGE_WAIT_TIMEOUT, { sender_id, message_id, chat_id });
    }
  }
  else {
    // Public, bot must be mentioned
    console.log('public chat');

    if (!_.some(KEYWORDS, x => _.toLower(message.text).includes(_.toLower(x)))) {
      // Bot not mentioned, ignore

      console.log("Bot not mentioned, let's not get this bread.");
      return;
    }
  }

  // We've gotten here, which means that the bot either had a message sent directly to it, or it has been mentioned publicly
  if (!SHOULD_WAIT) {
    if (match.found) {
      updateReminders(sender_id, setTimeout(() => { forwardMessage(params); }, match.wait));

      // Confirm to sender that reminder is set.
      let text = 'Reminder set for ' + match.num + ' ' + match.units + ' from now.';
      if (!REPLY && !FORWARD && !params.assembled) text = `No message specified. ${text}`;

      sendMessage({
        chat_id,
        reply_to_message_id: message_id,
        text,
      });
    }
    else {
      sendMessage({
        chat_id,
        reply_to_message_id: message_id,
        text: MATCH_NOT_FOUND_TEXT,
      });
    }
  }
});

http.createServer(app).listen(8081);
console.log('Listening on port 443.....');
