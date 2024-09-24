const getUserInfo = (req, res) => {
  if (!req.isAuthenticated()) {
    return res.status(401).json({ error: 'Unauthorized. Please login.' });
  }
  
  const user = req.user;  // User from OAuth2 session
  res.json({ username: user.username, discriminator: user.discriminator });
};

const destroySession = (req, res) => {
  req.logout();  // Passport logout (optional, depending on Passport usage)
  req.session.destroy((err) => {
    if (err) {
      console.error('Error destroying session:', err);
      return res.status(500).send('Logout failed.');
    }
    res.clearCookie('connect.sid');  // Clear the session cookie
    //TODO Check this out
    //res.redirect('/');               // Redirect user after logout
  });
}

  
module.exports = { getUserInfo, destroySession};