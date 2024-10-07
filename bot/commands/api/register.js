const {SlashCommandBuilder} = require('discord.js');
const axios = require('axios');

// Command to create a user
module.exports = {
    data: new SlashCommandBuilder()
        .setName('register')
        .setDescription('Register a new user in the database'),
    async execute(interaction) {
        const userId = interaction.user.id;
        const username = interaction.user.username;

        try {
            // Call Go backend to create a new user
            const response = await axios.post('http://localhost:8080/v1/api/users', {
                discord_user_id: userId,
            });

            if (response.status === 201) {
                await interaction.reply(`User ${username} registered successfully!`);
            } else {
                await interaction.reply(`Could not register user.`);
            }
        } catch (error) {
            console.error(error);
            await interaction.reply('An error occurred while registering the user.');
        }
    },
};