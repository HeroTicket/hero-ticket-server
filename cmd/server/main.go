package main

import (
	"github.com/heroticket/internal/cmd"
)

//	@title			Hero Ticket API
//	@version		1.0
//	@description	API for Hero Ticket DApp
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		api.heroticket.xyz
// @BasePath	/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Run()
}
