const express = require('express');
const router = express.Router();
const { getUserInfo } = require('../controllers/authController');

// Protected route to get user info
router.get('/user_info', getUserInfo);

module.exports = router;
