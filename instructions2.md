### Adding Authentication and Permissions to Your Taskify App

Now that your Taskify app is ready with the Tasks API, it’s time to add authentication and permissions to secure the application and control user access. This instruction will guide you through integrating **JWT-based authentication** and **Casbin for permissions management** in your existing Gin app.

---

### **Step 1: Install Required Packages**
Install the libraries needed for authentication and permissions:
```bash
go get github.com/golang-jwt/jwt/v5
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3  # Optional: DB-backed policies
go get golang.org/x/crypto/bcrypt        # For password hashing
```

---

### **Step 2: Set Up Authentication**

#### **a. Define User Model**
Extend your app to include a `User` model with fields for `Username`, `Password` (hashed), and `Role`:
- Fields: `ID`, `Username`, `Password`, `Role`.

#### **b. Create Authentication Handlers**
1. **Register Handler**:
   - Accept `username`, `password`, and `role` in the request body.
   - Hash the password using `bcrypt` and store it in the database.
   - Assign the user a role (e.g., `admin`, `editor`, `viewer`).

2. **Login Handler**:
   - Verify the username and password.
   - Generate a JWT token on successful authentication.

3. **Token Generation**:
   - Use `golang-jwt/jwt` to create JWT tokens.
   - Include `username` and `role` in the token claims.

4. **Middleware to Verify Tokens**:
   - Create a middleware function to validate JWT tokens for protected routes.
   - Extract user information from the token claims and attach it to the Gin context.

---

### **Step 3: Add Permissions with Casbin**

#### **a. Initialize Casbin**
1. Define your Casbin model (`model.conf`):
   ```ini
   [request_definition]
   r = sub, obj, act

   [policy_definition]
   p = sub, obj, act

   [role_definition]
   g = _, _

   [policy_effect]
   e = some(where (p.eft == allow))

   [matchers]
   m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
   ```

2. Define policies in a file or database (e.g., `policy.csv`):
   ```csv
   p, admin, /tasks, GET
   p, admin, /tasks, POST
   p, admin, /tasks, DELETE
   p, editor, /tasks, GET
   p, viewer, /tasks, GET
   g, user1, admin
   g, user2, editor
   ```

3. Load the model and policies:
   - File-based:
     ```go
     e, _ := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")
     ```
   - Database-backed:
     ```go
     adapter := gormadapter.NewAdapter("mysql", "user:password@tcp(127.0.0.1:3306)/dbname", true)
     e, _ := casbin.NewEnforcer("path/to/model.conf", adapter)
     ```

#### **b. Middleware for Permission Enforcement**
- Create middleware to check permissions using Casbin's `Enforce` method:
  ```go
  allowed, _ := e.Enforce(userRole, resource, action)
  ```
- Deny access if `allowed` is `false`.

---

### **Step 4: Update Your Routes**

1. **Public Routes**:
   - `/register`: User registration.
   - `/login`: User login and token generation.

2. **Protected Routes**:
   - Use the JWT middleware to validate tokens.
   - Use the Casbin middleware to enforce permissions.

3. **Example Route Registration**:
   ```go
   r.POST("/register", registerHandler)
   r.POST("/login", loginHandler)

   taskRoutes := r.Group("/tasks")
   taskRoutes.Use(jwtMiddleware)
   taskRoutes.Use(casbinMiddleware)
   {
       taskRoutes.GET("", getTasksHandler)
       taskRoutes.POST("", createTaskHandler)
       taskRoutes.DELETE("/:id", deleteTaskHandler)
   }
   ```

---

### **Step 5: Testing**

#### **a. Register and Login**
1. Register a user with a specific role (e.g., `admin`).
2. Log in to obtain a JWT token.

#### **b. Test Protected Routes**
1. Access protected routes with the token.
2. Verify that permissions are enforced correctly:
   - Admins can perform all actions.
   - Editors can only create and view tasks.
   - Viewers can only view tasks.

---

### **Step 6: Best Practices**
1. **Secure JWT Tokens**:
   - Use a strong secret key and environment variables to store it.
   - Set a reasonable expiration time for tokens.

2. **Role Management**:
   - Use Casbin to dynamically update roles and permissions as your app scales.
   - Store roles and permissions in a database for easier management.

3. **Middleware Performance**:
   - Cache frequently used Casbin policies to reduce latency.
   - Load Casbin policies and models during app initialization.

---

### **Final Workflow**

1. Users register and log in to obtain a JWT token.
2. JWT middleware validates the token for protected routes.
3. Casbin middleware enforces permissions based on the user's role and policies.