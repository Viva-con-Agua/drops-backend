db = db.getSiblingDB("drops")
db.user.createIndex({ email: 1 }, { unique: true })

