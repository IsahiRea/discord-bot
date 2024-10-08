const {SlashCommandBuilder, PermissionFlagsBits} = require('discord.js');

// Permissions
module.exports = {
    data: new SlashCommandBuilder()
        .setName('kick')
        .setDescription('Select a member and kick them.')
        .addUserOption(option =>
            option
                .setName('target')
                .setDescription('The member to kick')
                .setRequired(true))
        .addStringOption(option =>
            option
                .setName('reason')
                .setDescription('The reason for banning'))
        .setDefaultMemberPermissions(PermissionFlagsBits.KickMembers),
        //.setDMPermission(false)
    async execute(interaction) {
        const target = interaction.options.getUser('target');
        const reason = interaction.options.getString('reason') ?? 'No reason provided';

        await interaction.reply(`Kicking ${target.username} for reason: ${reason}`);
        await interaction.guild.members.ban(target);
    },
};