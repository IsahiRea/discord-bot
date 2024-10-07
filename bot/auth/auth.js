const jwt = require('jsonwebtoken');

//FIXME: Refactor code to work with refresh tokens

// Function to generate a JWT token
function generateToken(userId) {
    // Define the payload
    const payload = { userId: userId };

    // Sign the token with a secret key (stored in environment variable)
    const token = jwt.sign(payload, process.env.JWT_SECRET, { expiresIn: '1h' }); // Token expires in 1 hour

    return token;
}

// Function to generate access and refresh tokens
async function authenticateBot() {
    try {
        const response = await axios.post('http://localhost:8080/login', {
            clientId: process.env.CLIENT_ID,  // Example: use client credentials
            clientSecret: process.env.CLIENT_SECRET
        });

        accessToken = response.data.accessToken;
        refreshToken = response.data.refreshToken;

        console.log('Access Token:', accessToken);
        console.log('Refresh Token:', refreshToken);
    } catch (error) {
        console.error('Failed to authenticate the bot:', error);
    }
}

// Function to refresh the access token using the refresh token
async function refreshAccessToken() {
    try {
        const response = await axios.post('http://localhost:8080/refresh', {
            refreshToken: refreshToken
        });

        accessToken = response.data.accessToken;
        console.log('New Access Token:', accessToken);
    } catch (error) {
        console.error('Failed to refresh access token:', error);
        // Handle re-authentication if refresh token is invalid or expired
    }
}

module.exports = generateToken;