const express = require('express');
const session = require('express-session');
const passport = require('passport');
require('./config/passportConfig');
const apiRoutes = require('./routes/api');
const connectDB = require('./config/connectDB');
const MongoStore = require('connect-mongo');
const { expiration } = require('./middleware/expiration');
require('dotenv').config()

const app = express();

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

connectDB();

// Session setup (using connect-mongo for session store)
app.use(session({
  secret: process.env.SESSION_SECRET,  // Keep this secret safe (e.g., in .env)
  resave: false,                       // Prevent session resaving if unmodified
  saveUninitialized: false,            // Prevent saving uninitialized sessions
  store: MongoStore.create({
    mongoUrl: process.env.MONGO_URI,   // MongoDB connection URI
    collectionName: 'sessions'         // Name of the collection where sessions will be stored
  }),
  cookie: {
    httpOnly: true,                      // Prevents JavaScript access to cookies
    secure: process.env.NODE_ENV === 'production',  // Send cookies only over HTTPS in production
    sameSite: 'lax',  
    maxAge: 1000 * 60 * 60 * 24 * 7,   // 1-week session expiration
  }
}));

app.use(expiration);


// Passport setup
app.use(passport.initialize());
app.use(passport.session());

// API routes
app.use('/api', apiRoutes);

PORT = process.env.PORT || 5000
app.listen(PORT, () => {
    console.log(`Backend server running at http://localhost:${PORT}`);
  });