const {SlashCommandBuilder} = require('discord.js')
const axios = require('axios')

module.exports = {
    data: new SlashCommandBuilder()
        .setName('imagine')
        .setDescription('Create an AI Image'),
    async execute(interaction) {

        message = interaction.message
        try {
            const response = await axios.post('http://localhost:3000/api/images', {
                message: message 
            });

            // TODO Display the image correctly
            story = response.data.image
            await interaction.reply(image)   
        } catch (error) {
            console.log(error)
            await interaction.reply('An error occured while creating image.')
        }
    }
}