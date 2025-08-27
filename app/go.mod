module main

replace initialize => ./../packages/system/init/module

replace note => ./../packages/services/note/module

replace drive => ./../packages/services/drive/module

go 1.25

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lemoras/goutils/api v1.0.3 // indirect
	github.com/lemoras/goutils/db v1.0.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

require (
	drive v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	initialize v0.0.0-00010101000000-000000000000
	note v0.0.0-00010101000000-000000000000
)
