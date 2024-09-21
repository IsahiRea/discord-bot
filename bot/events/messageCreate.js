const { Collection, Events } = require('discord.js');

module.exports = {
	name: Events.MessageCreate,
    async execute(message) {
        //TODO Work on exectuing API endpoints

        if (message.author.bot) return;

        const args = message.content.slice(1).trim().split(/ +/);
        const commandName = args.shift().toLowerCase();
    
        if (!client.commands.has(commandName)) return;

        try {
            await command.execute(message, args);
        } catch (error) {
            console.error(error);
            message.channel.send('There was an error executing that command.');
        }
    }
}