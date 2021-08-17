package price

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/karlmutch/aws-ec2-price/pkg/price/version"
)

var (
	URL         = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
	URL_VERSION = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/index.json"

	HOURLY_TERM_CODE = "JRTCKXETXF"
	RATE_CODE        = "6YS6EN2CT7"

	REQUIRED_TERM           = "OnDemand"
	REQUIRED_TENANCY        = "Shared"
	REQUIRED_PRODUCT_FAMILY = "Compute Instance"
	REQUIRED_OS             = "Linux"
	REQUIRED_LICENSE_MODEL  = "No License required"
	REQUIRED_USAGE          = "BoxUsage:*"
	REQUIRED_PREINSTALLEDSW = "NA"

	CACHED_PRICING = CachedEc2Pricing{
		infos: map[string]*ec2PricingInfo{},
	}
)

type ec2Pricing struct {
	Products map[string]Ec2Product `json:products`
	Terms    map[string]map[string]map[string]struct {
		PriceDimensions map[string]struct {
			PricePerUnit struct {
				USD string `json:USD`
			} `json:pricePerUnit`
		} `json:priceDimensions`
	} `json:terms`
}

func (ec *ec2Pricing) GetInstances(region string) ([]*Instance, error) {
	var instances []*Instance
	for _, product := range ec.Products {
		if product.isValid() == false {
			continue
		}

		if product.isValidRegion(region) == false {
			continue
		}

		h := fmt.Sprintf("%s.%s", product.Sku, HOURLY_TERM_CODE)
		r := fmt.Sprintf("%s.%s.%s", product.Sku, HOURLY_TERM_CODE, RATE_CODE)

		usd := ec.Terms[REQUIRED_TERM][product.Sku][h].PriceDimensions[r].PricePerUnit.USD

		price, err := strconv.ParseFloat(usd, 64)
		if err != nil {
			return nil, errors.New("usd could not be parsed.")
		}

		instances = append(instances, &Instance{
			Region: region,
			Type:   product.InstanceType(),
			Price:  price,
		})
	}

	return instances, nil
}

func (ec *ec2Pricing) GetInstance(region string, instanceType string) (*Instance, error) {
	instances, err := ec.GetInstances(region)
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.Type != instanceType {
			continue
		}

		return instance, nil
	}

	return nil, errors.New("there is no matched instance.")
}

type Ec2Product struct {
	Sku           string `json:sku`
	ProductFamily string `json:productFamily`
	Attributes    struct {
		Location        string `json:location`
		InstanceType    string `json:instanceType`
		Tenancy         string `json:tenancy`
		OperatingSystem string `json:operatingSystem`
		LicenseModel    string `json:licenseModel`
		UsageType       string `json:usagetype`
		PreInstalledSw  string `json:preInstalledSw`
	}
}

func (ep *Ec2Product) InstanceType() string {
	return ep.Attributes.InstanceType
}

func (ep *Ec2Product) isValidRegion(region string) bool {
	if r, ok := REGIONS[region]; ok == true {
		return ep.Attributes.Location == r
	}

	return false
}

func (ep *Ec2Product) isValid() bool {
	if ep.ProductFamily != REQUIRED_PRODUCT_FAMILY {
		return false
	}

	if ep.Attributes.OperatingSystem != REQUIRED_OS {
		return false
	}

	if ep.Attributes.LicenseModel != REQUIRED_LICENSE_MODEL {
		return false
	}

	if ep.Attributes.Tenancy != REQUIRED_TENANCY {
		return false
	}

	if ep.Attributes.PreInstalledSw != REQUIRED_PREINSTALLEDSW {
		return false
	}

	matched, err := regexp.MatchString(REQUIRED_USAGE, ep.Attributes.UsageType)
	if err != nil || matched == false {
		return false
	}

	return true

}

type CachedEc2Pricing struct {
	infos map[string]*ec2PricingInfo
}

type ec2PricingInfo struct {
	pricing       *ec2Pricing
	lastCheckTime time.Time
}

func (c *ec2PricingInfo) Update(pricing *ec2Pricing) {
	c.pricing = pricing
}

func (c *CachedEc2Pricing) isExpired(region string) bool {
	if val, ok := c.infos[region]; ok {
		return time.Since(val.lastCheckTime) > time.Duration(24*time.Hour)
	}
	return true
}

func (c *CachedEc2Pricing) update(region string, pricing *ec2Pricing) {
	if val, ok := c.infos[region]; ok {
		val.Update(pricing)
	} else {
		e := &ec2PricingInfo{
			pricing:       pricing,
			lastCheckTime: time.Now(),
		}
		c.infos[region] = e
	}
}

func NewPricing(region string) (*ec2Pricing, error) {
	if CACHED_PRICING.isExpired(region) == false {
		return CACHED_PRICING.infos[region].pricing, nil
	}

	client := &http.Client{}

	r, err := client.Get(URL_VERSION)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	versions := &version.Version{}
	if err := json.NewDecoder(r.Body).Decode(versions); err != nil {
		return nil, err
	}

	fmt.Println(versions.CurrentVersion)
	url := fmt.Sprintf("https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/%s/%s/index.json", versions.CurrentVersion, region)
	fmt.Println(url)

	r, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	pricing := &ec2Pricing{}
	if err := json.NewDecoder(r.Body).Decode(pricing); err != nil {
		return nil, err
	}

	CACHED_PRICING.update(region, pricing)
	return pricing, err
}
