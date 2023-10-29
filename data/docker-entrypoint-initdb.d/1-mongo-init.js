var user = {
  user: "mongo-user",
  pwd: "user-password",
  roles: [
    {
      role: "dbOwner",
      db: "todofuken-db"
    }
  ]
};

db.createUser(user);
db.createCollection('todofuken');
