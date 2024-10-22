const axios = require('axios');
require("dotenv").config()

// Function to generate access and refresh tokens
async function authenticateBot(discordID) {
    try {
        const response = await axios.post('http://localhost:3000/api/login', {
            discord_id: discordID,
            client_id: process.env.clientID,
        });

        const authData = response.data;
        const refreshToken = authData.refresh_token;
        console.log('Refresh Token:', refreshToken);
        return refreshToken;
    } catch (error) {
        console.error('Failed to authenticate the bot:', error);
    }
}

// Function to refresh the access token using the refresh token
async function refreshAccessToken(refreshToken) {
    try {
        const response = await axios.get('http://localhost:3000/api/refresh', {
            headers:{
                'Authorization': `Bearer ${refreshToken}`
            }
        });

        const authData = response.data;
        const accessToken = authData.access_token;
        console.log('New Access Token:', accessToken);
        return accessToken;
    } catch (error) {
        console.error('Failed to refresh access token:', error);
    }

    return accessToken;
}

module.exports = {authenticateBot, refreshAccessToken};