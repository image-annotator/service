# Go Restful API Boilerplate

[![GoDoc Badge]][GoDoc] [![GoReportCard Badge]][GoReportCard]

Easily extendible RESTful API boilerplate aiming to follow idiomatic go and best practice.

The goal of this boiler is to have a solid and structured foundation to build upon on.

### Features
The following feature set is a minimal selection of typical Web API requirements:

- Configuration using [viper](https://github.com/spf13/viper)
- CLI features using [cobra](https://github.com/spf13/cobra)
- PostgreSQL support including migrations using [go-pg](https://github.com/go-pg/pg)
- Structured logging with [Logrus](https://github.com/sirupsen/logrus)
- Routing with [chi router](https://github.com/go-chi/chi) and middleware
- JWT Authentication using [jwt-go](https://github.com/dgrijalva/jwt-go) in combination with passwordless email authentication (could be easily extended to use passwords instead)
- Request data validation using [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- HTML emails with [gomail](https://github.com/go-gomail/gomail)

### Start Application
- Clone this repository
- Create a postgres database and set environment variable *DATABASE_URL* accordingly if not using same as default
- Build the application: ```go build``` to create ```go-base``` binary
- Initialize the database and run all migrations found in ./database/migrate with: ```go-base migrate```
- Run the application: ```go-base serve```

### API Routes

For passwordless login following routes are available:

Path | Method | Required JSON | Header | Description
---|---|---|---|---
/auth/login | POST | email | | the email you want to login with (see below)
/auth/token | POST | token | | the token you received via email (or printed to stdout if smtp not set)
/auth/refresh | POST | | Authorization: "Bearer refresh_token" | refresh JWTs
/auth/logout | POST | | Authorizaiton: "Bearer refresh_token" | logout from this device

Besides /auth/* the API provides to main routes /api/* and /admin/* to distinguish between application and administration features. The latter requires to be logged in as administrator by providing the respective JWT in Authorization Header.

Check [routes.md](routes.md) file for an overview of the provided API routes.


### Client API Access and CORS
The server is configured to serve a Progressive Web App (PWA) client from it's "public" folder. this is where you put the contents of your client's build "dist" folder into. In this case enabling CORS is not required, because the client is served from the same host as the api.

If you want to access the api from a client that is serverd from a different host, including e.g. a development live reloading server, you must enable CORS on the server first by setting an environment variable ENABLE_CORS=true for the server to acceept api connections from clients serverd by other hosts.

#### Demo client application
For demonstration of the login and account management features this API serves a demo [Vue.js](https://vuejs.org) PWA. The client's source code can be found [here](https://gitlab.informatika.org/label-1-backend/base-vue).

If no valid email smtp settings are provided by environment variables, emails will be print to stdout showing the login token. Use one of the following bootstrapped users for login:
- admin@boot.io (has access to admin panel)
- user@boot.io

A deployed version can also be found on [Heroku](https://govue.herokuapp.com)

### Environment Variables

Name | Type | Default | Description
---|---|---|---
PORT | string | localhost:3000 | http address (accepts also port number only for heroku compability)  
LOG_LEVEL | string | debug | log level
LOG_TEXTLOGGING | bool | false | defaults to json logging
DB_NETWORK | string | tcp | database 'tcp' or 'unix' connection
DB_ADDR | string | localhost:5432 | database tcp address or unix socket
DB_USER | string | postgres | database user name
DB_PASSWORD | string | postgres | database user password
DB_DATABASE | string | gobase | database shema name
AUTH_LOGIN_URL | string | http://localhost:3000/login | client login url as sent in login token email
AUTH_LOGIN_TOKEN_LENGTH | int | 8 | length of login token
AUTH_LOGIN_TOKEN_EXPIRY | time.Duration | 11m | login token expiry
AUTH_JWT_SECRET | string | random | jwt sign and verify key - value "random" creates random 32 char secret at startup (and automatically invalidates existing tokens on app restarts, so during dev you might want to set a fixed value here)
AUTH_JWT_EXPIRY | time.Duration | 15m | jwt access token expiry
AUTH_JWT_REFRESH_EXPIRY | time.Duration | 1h | jwt refresh token expiry
EMAIL_SMTP_HOST | string || email smtp host (if set and connection can't be established then app exits)
EMAIL_SMTP_PORT | int || email smtp port
EMAIL_SMTP_USER | string || email smtp username
EMAIL_SMTP_PASSWORD | string || email smtp password
EMAIL_FROM_ADDRESS | string || from address used in sending emails
EMAIL_FROM_NAME | string || from name used in sending emails
ENABLE_CORS | bool | false | enable CORS requests

### Contributing

Any feedback and pull requests are welcome and highly appreciated. Please open an issue first if you intend to send in a larger pull request or want to add additional features.

[GoDoc]: https://godoc.org/gitlab.informatika.org/label-1-backend/base
[GoDoc Badge]: https://godoc.org/gitlab.informatika.org/label-1-backend/base?status.svg
[GoReportCard]: https://goreportcard.com/report/gitlab.informatika.org/label-1-backend/base
[GoReportCard Badge]: https://goreportcard.com/badge/gitlab.informatika.org/label-1-backend/base