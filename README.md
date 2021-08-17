# aws-ec2-price

Go package, and REST API for searching the prices of AWS EC2 instances

Installation

------------
	go get -u github.com/karlmutch/aws-ec2-price

This package requires some dependancies:
* [gin-gonic/gin](https://github.com/gin-gonic/gin)
* [urfave/cli](https://github.com/urfave/cli)

Usage in golang
---------------
```go
import "github.com/karlmutch/aws-ec2-price"

pricing, err := price.NewPricing("us-east-1")
if err != nil {
	doSomething()
}

instances, err := pricing.GetInstances("us-east-1")
if err != nil {
	doSomething()
}

println(instances)

instance, err := pricing.GetInstance("us-east-1", "c4.large")
if err != nil {
	doSomething()
}

println(instance)
```


Build
-----

make all

Run 
---
	aws-ec2-price --port=[PORT]

Default port is 8080.

Run with Docker
---------------
	docker run -p 8080:8080 -d karlmutch/aws-ec2-price

REST APIs
---------
* GET /ec2/regions/:region
* GET /ec2/regions/:region/instance_types/:instance_type

References
----------
* https://aws.amazon.com/ko/blogs/aws/new-aws-price-list-api/
