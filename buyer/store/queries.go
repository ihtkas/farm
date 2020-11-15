package store

// TODO: Evaluate otehr ways to make this externalise. Consider remodelling in updates to app in production, migration etc.,

const (
	orderTable     = "product_order"
	userOrderTable = "user_order"

	orderIDColumn           = "order_id"
	productIDColumn         = "product_id"
	userIDColumn            = "user_id"
	quantityColumn          = "quantity"
	locLatColumn            = "loc_lat"
	locLonColumn            = "loc_lon"
	productProtoColumn      = "product_proto"
	insertionTimeUUIDColumn = "insertion_time"

	createOrderCasandraQuery = "CREATE TABLE IF NOT EXISTS " + orderTable + " (" +
		orderIDColumn + " UUID PRIMARY KEY," +
		productIDColumn + " UUID," +
		quantityColumn + " int," +
		locLatColumn + " double," +
		locLonColumn + " double)"

	createUserOrderCasandraQuery = "CREATE TABLE IF NOT EXISTS " + userOrderTable + " (" +
		userIDColumn + " UUID," +
		insertionTimeUUIDColumn + " UUID," +
		orderIDColumn + " UUID," +
		productIDColumn + " UUID," +
		quantityColumn + " int," +
		locLatColumn + " double," +
		locLonColumn + " double," +
		"PRIMARY KEY ((" + userIDColumn + "), " + insertionTimeUUIDColumn + "))" +
		"WITH CLUSTERING ORDER BY (" + insertionTimeUUIDColumn + " DESC)"

	insertProductCassandraQuery = "INSERT INTO " + orderTable + "(" +
		orderIDColumn + "," +
		productIDColumn + "," +
		quantityColumn + "," +
		locLatColumn + "," +
		locLonColumn +
		") values (?, ?, ?, ?, ?)"

	insertUserOrderCassandraQuery = "INSERT INTO " + userOrderTable + "(" +
		userIDColumn + "," +
		insertionTimeUUIDColumn + "," +
		orderIDColumn + "," +
		productIDColumn + "," +
		quantityColumn + "," +
		locLatColumn + "," +
		locLonColumn + ") values (?, ?, ?, ?, ?, ?, ?)"

	selectOrderCassandraQuery = "SELECT " +
		productIDColumn + "," +
		quantityColumn + "," +
		locLatColumn + "," +
		locLonColumn +
		" FROM " + orderTable + " WHERE " + orderIDColumn + "=?"
	selectOrderByUserAfterTimeCassandraQuery = "SELECT " +
		insertionTimeUUIDColumn + "," +
		orderIDColumn + "," +
		productIDColumn + "," +
		quantityColumn + "," +
		locLatColumn + "," +
		locLonColumn +
		" FROM " + userOrderTable + " WHERE " + userIDColumn + "=? and " + insertionTimeUUIDColumn + ">? limit ?"

	selectOrderByUserCassandraQuery = "SELECT " +
		insertionTimeUUIDColumn + "," +
		orderIDColumn + "," +
		productIDColumn + "," +
		quantityColumn + "," +
		locLatColumn + "," +
		locLonColumn +
		" FROM " + userOrderTable + " WHERE " + userIDColumn + "=? limit ?"
)
