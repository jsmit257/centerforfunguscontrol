module github.com/jsmit257/centerforfunguscontrol

go 1.23.1

replace github.com/jsmit257/huautla => /home/johnny/dev/go/src/github.com/jsmit257/huautla

replace github.com/jsmit257/userservice => /home/johnny/dev/go/src/github.com/jsmit257/userservice

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/google/uuid v1.6.0
	github.com/jsmit257/huautla v0.0.0-20241111220754-08460da5815b
	github.com/jsmit257/userservice v0.0.0-20241119014602-e7422fe454fa
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
