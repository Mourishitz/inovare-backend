# Authentication and Authorization

This project has a very simple role-based authentication and authorization system to manage user access and permissions.
The roles defined in this system are:
- **Admin**: Has full access to all resources and can manage other users.
- **User**: Has limited access, and can only view and interact with certain resources.


## How it works
- Upon login, users are authenticated and assigned a role based on their credentials through a JWT token processed on `middleware/auth.go`
- This middleware can be called on any route to protect it and future logic can be used by checking the variable inside `c.userID` and `c.username`;
- The roles are defined inside `models/roles.go` (but they are not stored on database).
- The available roles are numbers, starting at 1 (Admin) and 2 (User). New roles can be added by adding a new constant on `models/roles.go`.
- The default token expiration time is 24 hours, but it can be changed on `utils/jwt.go`
