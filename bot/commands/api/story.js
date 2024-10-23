const {SlashCommandBuilder} = require('discord.js')
const axios = require('axios')

module.exports = {
    data: new SlashCommandBuilder()
        .setName('story')
        .setDescription('Build a story with the server'),
    async execute(interaction) {

        message = interaction.message
        try {
            const response = await axios.post('http://localhost:3000/api/stories', {
                message: message 
            });

            story = response.data.story
            await interaction.reply(story)   
        } catch (error) {
            console.log(error)
            await interaction.reply('An error occured while adding to the story.')
        }
    }
}