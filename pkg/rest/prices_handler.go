package rest

import (
	"github.com/karlmutch/aws-ec2-price/pkg/price"

	"github.com/gin-gonic/gin"
)

func getEc2PricesHandler(context *gin.Context) {
	region := context.Param("region")

	pricing, err := price.NewPricing(region)
	if err != nil {
		context.Error(err)
		errorHandler(context)
		return
	}

	instances, err := pricing.GetInstances(region)
	if err != nil {
		context.Error(err)
		errorHandler(context)
		return
	}

	context.JSON(200, instances)
}
