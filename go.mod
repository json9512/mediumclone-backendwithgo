module json9512/mediumclone-go

go 1.15

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/franela/goblin v0.0.0-20201006155558-6240afcb2eb7
	github.com/gin-gonic/gin v1.6.3
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/toorop/gin-logrus v0.0.0-20200831135515-d2ee50d38dae
	gorm.io/gorm v1.20.11 // indirect
	json9512/mediumclone-go/db v0.0.0-00010101000000-000000000000
	json9512/mediumclone-go/util v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	json9512/mediumclone-go/db => ./src/db
	json9512/mediumclone-go/util => ./src/util
)
