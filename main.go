package main

import "github.com/iamarpitzala/aca-reca-backend/cmd"

// @title           ACA RECA API
// @version         1.0
// @description     API documentation for ACA RECA Backend
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@acareca.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cmd.InitServer()
}
