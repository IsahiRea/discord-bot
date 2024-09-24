const express = require('express');
const router = express.Router();
const { getUserInfo, destroySession } = require('../controllers/authController');

// Protected route to get user info
router.get('/user_info', getUserInfo);
router.get('/logout', destroySession);

module.exports = router;
