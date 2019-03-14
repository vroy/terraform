package test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceAsSingle(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_as_single" "foo" {
	resource_as_block {
		foo = "as block a"
	}
	resource_as_attr = {
		foo = "as attr a"
	}
	primitive = "primitive a"
}
				`),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("state after initial create:\n", s.String())
						return nil
					},
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_block.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_block.0.foo", "as block a"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_attr.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_attr.0.foo", "as attr a"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "primitive.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "primitive.0", "primitive a"),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_as_single" "foo" {
	resource_as_block {
		foo = "as block b"
	}
	resource_as_attr = {
			foo = "as attr b"
	}
	primitive = "primitive b"
}
				`),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("state after update:\n", s.String())
						return nil
					},
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_block.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_block.0.foo", "as block b"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_attr.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_attr.0.foo", "as attr b"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "primitive.#", "1"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "primitive.0", "primitive b"),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_as_single" "foo" {
}
				`),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("state after everything unset:\n", s.String())
						return nil
					},
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_block.#", "0"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "resource_as_attr.#", "0"),
					resource.TestCheckResourceAttr("test_resource_as_single.foo", "primitive.#", "0"),
				),
			},
		},
	})
}
