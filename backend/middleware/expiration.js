// Middleware to enforce absolute session expiration (1 week)
function expiration (req, res, next) {
    if (req.session) {
        const now = Date.now();
        const maxSessionAge = 1000 * 60 * 60 * 24 * 7;  // 1 week
    
        // Check if session was created or renewed more than 1 week ago
        if (!req.session.createdAt) {
          req.session.createdAt = now;  // Set session creation timestamp if not set
        } else if (now - req.session.createdAt > maxSessionAge) {
          req.session.destroy((err) => {
            if (err) return next(err);
            res.clearCookie('connect.sid');  // Clear session cookie
            return res.status(401).send('Session expired, please log in again.');
          });
        }
      }
      next();
}

module.exports = expiration;
