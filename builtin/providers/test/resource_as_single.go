package test

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func testResourceAsSingle() *schema.Resource {
	return &schema.Resource{
		Create: testResourceAsSingleCreate,
		Read:   testResourceAsSingleRead,
		Delete: testResourceAsSingleDelete,
		Update: testResourceAsSingleUpdate,

		Schema: map[string]*schema.Schema{
			"resource_as_block": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				AsSingle: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foo": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"resource_as_attr": {
				Type:       schema.TypeList,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				MaxItems:   1,
				AsSingle:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foo": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"primitive": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				AsSingle: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func testResourceAsSingleCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId("placeholder")
	return testResourceAsSingleRead(d, meta)
}

func testResourceAsSingleRead(d *schema.ResourceData, meta interface{}) error {
	for _, k := range []string{"resource_as_block", "resource_as_attr", "primitive"} {
		v := d.Get(k)
		if v == nil {
			continue
		}
		if l, ok := v.([]interface{}); !ok {
			return fmt.Errorf("%s should appear as []interface{}, not %T", k, l)
		} else {
			for i, item := range l {
				switch k {
				case "primitive":
					if _, ok := item.(string); item != nil && !ok {
						return fmt.Errorf("%s[%d] should appear as string, not %T", k, i, item)
					}
				default:
					if _, ok := item.(map[string]interface{}); item != nil && !ok {
						return fmt.Errorf("%s[%d] should appear as map[string]interface{}, not %T", k, i, item)
					}
				}
			}
		}
	}
	return nil
}

func testResourceAsSingleUpdate(d *schema.ResourceData, meta interface{}) error {
	return testResourceAsSingleRead(d, meta)
}

func testResourceAsSingleDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
