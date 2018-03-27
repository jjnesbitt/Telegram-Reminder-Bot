var express = require('express');
var app = express();
var bodyParser = require('body-parser');
const axios = require('axios')

const BOT_TOKEN = '459017535:AAFl8jUbmKWHYGA2aAtE7tlHBRcStXjFoG0';


//In milliseconds
const time_units = {
    minute: 60000,
    hour: 3600000,
    day: 86400000,
    week: 604800000,
    month: 2419200000,
    year: 29030400000
}

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({
    extended: true
})); // for parsing application/x-www-form-urlencoded


//Function called to forward message, usually set on a timer.
function forwardMessage(chat_id, from_chat_id, message_id){
    console.log('yoooo');
}

//This is the route the API will call
app.post('/new-message', function(req, res) {
    const {message} = req.body

    var match = null;

    if (!message) {
        return res.end()
    }
    else{
        var matched = false;

        for (var key in time_units){
            var regex = new RegExp("(\\d{1,10}.?\\d{0,10}) " + key);
            if ((match = regex.exec(message)) != null){
                //Means valid number has been matched

                var num = match[1];
                match.push(key);
            }
        }

        if (!match){
            return res.end();
        }
    }

    var timeToWait = match[1]*time_units[match[match.length-1]];

    setTimeout(forwardMessage(message.from.id, message.chat.id, message.id), timeToWait);
    // Respond by hitting the telegram bot API and responding to message.
    axios.post('https://api.telegram.org/bot' + BOT_TOKEN + '/sendMessage', {
        chat_id: message.chat.id,
        text: 'Reminder set for ' + match[1] + ' ' + match[match.length - 1] + 's' + 'from now.',
        reply_to_message: message
    })
        .then(response => {
            // We get here if the message was successfully posted
            console.log('Message posted')
            res.end('ok')
        })
        .catch(err => {
            // ...and here if it was not
            console.log('Error :', err)
            res.end('Error :' + err)
        })

});

// Finally, start our server
app.listen(3000, function() {
    console.log('Telegram app listening on port 3000!');
});
