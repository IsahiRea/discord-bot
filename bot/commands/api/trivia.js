const {SlashCommandBuilder, bold, StringSelectMenuBuilder} = require('discord.js')
const axios = require('axios')

module.exports = {
    data: new SlashCommandBuilder()
        .setName('trivia')
        .setDescription('Get a random trivia question'),
    async execute(interaction) {

        try {
            const response = await axios.get('http://localhost:3000/api/trivia')
            

            const select =new StringSelectMenuBuilder()
            .setCustomId('trivia_question')
            .a

            const triviaData = response.data
            await interaction.reply(`Question: ${triviaData.question}`)


        } catch (error) {
           console.log(error) 
           await interaction.reply('An error occured while obtaining trivia question.')
        }
    }
}