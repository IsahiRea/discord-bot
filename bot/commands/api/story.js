const {SlashCommandBuilder} = require('discord.js')
const axios = require('axios')

module.exports = {
    data: new SlashCommandBuilder()
        .setName('story')
        .setDescription('Build a story with the server'),
    async execute(interaction) {

    }
}