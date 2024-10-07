const {SlashCommandBuilder} = require('discord.js');
const {generateToken} = require("../../auth/auth.js")
const axios = require('axios');
require('dotenv').config();

//FIXME: Refactor code to work with Tokens

// Command to get a user's information
module.exports = {
    data: new SlashCommandBuilder()
        .setName('userInfo')
        .setDescription('Provides information about the user'),
    async execute(interaction) {

        //Switch with access token
        const userId = interaction.user.id;
        const token = generateToken(userId);




        try {
            // Call Go backend to fetch user data
            const response = await axios.get(`http://localhost:8080/v1/api/users/${userId}`,{
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            const userData = response.data;
            await interaction.reply(`User: ${userData.username}`);
        } catch (error) {

            // Catch the expired access token



            console.error(error);
            await interaction.reply('Could not fetch user data.');
        }
    },
};