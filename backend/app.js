const express = require('express');
const session = require('express-session');
const passport = require('passport');
require('./config/passportConfig');
const apiRoutes = require('./routes/api')
require('dotenv').config()

const app = express();

// Middleware
app.use(express.json());
app.use(session({ secret: `${process.env.SECRET}`, resave: false, saveUninitialized: false }));
app.use(passport.initialize());
app.use(passport.session());

// API routes
app.use('/api', apiRoutes);

PORT = process.env.PORT || 5000
app.listen(PORT, () => {
    console.log(`Backend server running at http://localhost:${PORT}`);
  });