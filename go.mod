module github.com/MicaelaJofre/gocourse_enrollment

go 1.24.4

require (
	github.com/MicaelaJofre/gocourse_domain v0.0.2
	github.com/MicaelaJofre/gocourse_meta v0.0.1
	gorm.io/gorm v1.30.0
)

require github.com/joho/godotenv v1.5.1

require (
	github.com/MicaelaJofre/go_course_sdk v0.0.1
	github.com/MicaelaJofre/go_lib_response v0.0.1
)

require github.com/ncostamagna/go_http_client v0.0.3 // indirect

require (
	github.com/go-kit/kit v0.12.0
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.27.0 // indirect
	gorm.io/driver/mysql v1.6.0
)

replace github.com/MicaelaJofre/gocourse_domain => ../gocourse_domain

replace github.com/MicaelaJofre/gocourse_meta => ../gocourse_meta
