//FIXME Implement Oauth / passport

const passport = require('passport');
const DiscordStrategy = require('passport-discord').Strategy;
const User = require('../models/User');
require('dotenv').config();  // Load environment variables

// Configure Passport to use the Discord strategy
passport.use(new DiscordStrategy({
    clientID: process.env.DISCORD_CLIENT_ID,
    clientSecret: process.env.DISCORD_CLIENT_SECRET,
    callbackURL: process.env.DISCORD_REDIRECT_URI,
    scope: ['identify', 'guilds']
},
async (accessToken, refreshToken, profile, done) => {
    try {
        // Check if the user already exists in the database
        let user = await User.findOne({ discordId: profile.id });

        if (!user) {
            // If the user doesn't exist, create a new record
            user = new User({
                discordId: profile.id,
                username: profile.username,
                discriminator: profile.discriminator,
                avatar: profile.avatar,
                accessToken: accessToken,
                refreshToken: refreshToken
            });
            await user.save();
        }

        return done(null, user);  // Pass user to Passport
    } catch (err) {
        console.error('Error finding or saving user:', err);
        return done(err, null);
    }
}));

// Serialize user into the session (store user ID in session)
passport.serializeUser((user, done) => {
    done(null, user.id);
});

// Deserialize user from the session (fetch the complete user info from session storage)
passport.deserializeUser(async (id, done) => {
    try {
        const user = await User.findById(id);
        done(null, user);
    } catch (err) {
        done(err, null);
    }
});
