const {SlashCommandBuilder} = require('discord.js');
const axios = require('axios');

// Command to get a user's information
module.exports = {
    data: new SlashCommandBuilder()
        .setName('userInfo')
        .setDescription('Provides information about the user'),
    async execute(interaction) {
        const userId = interaction.user.id;

        try {
            // Call Go backend to fetch user data
            const response = await axios.get(`http://localhost:8080/v1/api/users/${userId}`);
            const userData = response.data;
            
            await interaction.reply(`User: ${userData.username}`);
        } catch (error) {
            console.error(error);
            await interaction.reply('Could not fetch user data.');
        }
    },
};