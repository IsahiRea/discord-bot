const {SlashCommandBuilder} = require('discord.js');
const {authenticateBot, refreshAccessToken} = require("../../auth/auth.js")
const axios = require('axios');
require('dotenv').config();

// Command to get a user's information
module.exports = {
    data: new SlashCommandBuilder()
        .setName('userinfo')
        .setDescription('Provides information about the user'),
    async execute(interaction) {

        const userId = interaction.user.id;

        try {

            const refreshToken = await authenticateBot(userId);
            const accessToken = await refreshAccessToken(refreshToken);

            // Call Go backend to fetch user data
            const response = await axios.get(`http://localhost:3000/api/users/${userId}`,{
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                }
            });

            const userData = response.data;
            await interaction.reply(`UserID: ${userData.DiscordUserID}`);
        } catch (error) {

            // Catch the expired access token
            console.error(error);
            await interaction.reply('Could not fetch user data. Please try again');
        }
    },
};