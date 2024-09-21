const axios = require('axios');
require('discord.js');

//FIXME 
module.exports = {
    name: 'getuserinfo',
    description: 'Fetch user info from the backend using OAuth2 authentication.',
    async execute(message, args) {
        try {
            const response = await axios.get('http://localhost:5000/api/user_info', {
                withCredentials: true  // Sends session cookies for authentication
            });
            const userInfo = response.data;
            message.channel.send(`Username: ${userInfo.username}#${userInfo.discriminator}`);
        } catch (error) {
            message.channel.send('Please log in first! Use /login to authenticate.');
        }
    }
};
