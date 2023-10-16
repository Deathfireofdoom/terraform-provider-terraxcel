package main

import (
	"context"

	"terraform-provider-terraxcel/terraxcel"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), terraxcel.New, providerserver.ServeOpts{
		Address: "deathfirefodoom.com/edu/terraxcel",
	})
}
