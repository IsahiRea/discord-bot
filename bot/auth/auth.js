require("dotenv").config()

// Function to generate access and refresh tokens
async function authenticateBot(discordID) {
    try {
        const response = await axios.post('http://localhost:8080/login', {
            discord_id: discordID,
            client_id: process.env.clientID,
        });

        refreshToken = response.data.refresh_token;
        console.log('Refresh Token:', refreshToken);
    } catch (error) {
        console.error('Failed to authenticate the bot:', error);
    }

    return refreshToken;
}

// Function to refresh the access token using the refresh token
async function refreshAccessToken(refreshToken) {
    try {
        const response = await axios.post('http://localhost:8080/refresh', {
            refreshToken: refreshToken
        });

        accessToken = response.data.access_token;
        console.log('New Access Token:', accessToken);
    } catch (error) {
        console.error('Failed to refresh access token:', error);
        // TODO Handle re-authentication if refresh token is invalid or expired
    }

    return accessToken;
}

module.exports = {authenticateBot, refreshAccessToken};