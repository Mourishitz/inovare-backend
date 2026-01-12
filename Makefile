migrate: 
	atlas schema apply --env gorm -u "postgres://username:password@localhost:5432/inovare?sslmode=disable"
