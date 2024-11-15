const {SlashCommandBuilder} = require('discord.js')
const axios = require('axios')
const {authenticateBot, refreshAccessToken} = require("../../auth/auth.js")

module.exports = {
    data: new SlashCommandBuilder()
        .setName('imagine')
        .setDescription('Create an AI Image'),
    async execute(interaction) {

        const message = interaction.message
        try {

            // TODO Grab response and display image on discord
            const response = await axios.post('http://localhost:3000/api/images', {
                message: message 
            });

            const image = response.data.image
            await interaction.reply(image)   
        } catch (error) {
            console.log(error)
            await interaction.reply('An error occured while creating image.')
        }
    }
}