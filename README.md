# Telegram Reminder Bot

This is a bot for telegram that will remind users of things, by forwarding messages back to user after the specified time.

### How to use
- Add it to a chat by the username: `@OGReminderBot`
- Reply to a message and invoke the command `\remindme <wait-time>`. The bot must be invoked by a command if it's in a public chat or it will not respond.
- Forward/Send a message directly to the bot, and specified the wait time after it replies to you.


Reminder times can be of the format `<integer>` `[seconds, minutes, hours, days, weeks, months, years]`


### TODO
- [x] Add storing reminders in database for reminder listing, and early deletion
- [x] Add command to cancel reminders as mentioned above
- [ ] Fix code organization to remove so many global variables

<br/>
<br/>

