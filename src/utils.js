import axios from 'axios';
import _ from 'lodash';

import { BOT_TOKEN, BOT_USERNAME } from './constants';

// Object that stores the timout ID's for forwarding messages, by user id
const currentReminders = {};

export const cancelReminders = sender_id => {
  const timeouts = currentReminders[sender_id];
  if (timeouts !== undefined) {
    timeouts.forEach(val => clearTimeout(val));
    delete currentReminders[sender_id];
  }
};

export const updateReminders = (sender_id, timeout_id) => {
  if (currentReminders[sender_id] !== undefined) {
    currentReminders[sender_id] = [
      ...currentReminders[sender_id],
      timeout_id,
    ];
  }
  else currentReminders[sender_id] = [timeout_id];
};

export const forwardMessage = async params => {
  try {
    const res = await axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/forwardMessage', params);
    console.log('Message forwarded');
  } catch (err) {
    console.log('Error:', err.keys());
  }
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

export const checkCommands = message => {
  const matchCommandRegex = /^\/(\w+)/;
  const match = matchCommandRegex.exec(message.text);

  console.log(match);
  if (!match) return false;

  const command = match[1];
  if (command == 'start') return true;
  else if (command == 'cancel') {
    cancelReminders(message.from.id);
    sendMessage({
      chat_id: message.chat.id,
      reply_to_message_id: message.message_id,
      text: REMINDERS_CLEARED_TEXT,
    });
    return true;
  }
  else if (command == 'reminders') {
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

export const checkIfFromTelegram = req => {
  const REAL_IP_HEADER = 'x-real-ip';
  const ipToCheck = '149.154.167';

  const ipArray = req.headers[REAL_IP_HEADER].split('.');
  if (_.dropRight(ipArray).join('.') != ipToCheck) return false;
  if (parseInt(ipArray[3]) < 197 || parseInt(ipArray[3]) > 233) return false;

  return true;
};
