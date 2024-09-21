const getUserInfo = (req, res) => {
  if (!req.isAuthenticated()) {
    return res.status(401).json({ error: 'Unauthorized. Please login.' });
  }
  
  const user = req.user;  // User from OAuth2 session
  res.json({ username: user.username, discriminator: user.discriminator });
};

  
  module.exports = getUserInfo;