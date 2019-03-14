package plugin

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestFixupAsSingleInstanceStateInOut(t *testing.T) {
	tests := map[string]struct {
		AttrsIn  map[string]string
		AttrsOut map[string]string
		Schema   map[string]*schema.Schema
	}{
		"empty": {
			nil,
			nil,
			nil,
		},
		"simple": {
			map[string]string{
				"a": "a value",
			},
			map[string]string{
				"a": "a value",
			},
			map[string]*schema.Schema{
				"a": {Type: schema.TypeString, Optional: true},
			},
		},
		"normal list of primitive, empty": {
			map[string]string{
				"a.%": "0",
			},
			map[string]string{
				"a.%": "0",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"normal list of primitive, single": {
			map[string]string{
				"a.%": "1",
				"a.0": "hello",
			},
			map[string]string{
				"a.%": "1",
				"a.0": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"AsSingle list of primitive": {
			map[string]string{
				"a": "hello",
			},
			map[string]string{
				"a.#": "1",
				"a.0": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					AsSingle: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"AsSingle list of resource": {
			map[string]string{
				"a.b": "hello",
			},
			map[string]string{
				"a.#":   "1",
				"a.0.b": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					AsSingle: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"b": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
		"AsSingle list of resource with nested primitive list": {
			map[string]string{
				"a.b.#": "1",
				"a.b.0": "hello",
			},
			map[string]string{
				"a.#":     "1",
				"a.0.b.#": "1",
				"a.0.b.0": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					AsSingle: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"b": {
								Type:     schema.TypeList,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
		"AsSingle list of resource with nested AsSingle primitive list": {
			map[string]string{
				"a.b": "hello",
			},
			map[string]string{
				"a.#":     "1",
				"a.0.b.#": "1",
				"a.0.b.0": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					AsSingle: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"b": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								AsSingle: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
		"AsSingle list of resource with nested AsSingle resource list": {
			map[string]string{
				"a.b.c": "hello",
			},
			map[string]string{
				"a.#":       "1",
				"a.0.b.#":   "1",
				"a.0.b.0.c": "hello",
			},
			map[string]*schema.Schema{
				"a": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					AsSingle: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"b": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								AsSingle: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"c": {
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	copyMap := func(m map[string]string) map[string]string {
		if m == nil {
			return nil
		}
		ret := make(map[string]string, len(m))
		for k, v := range m {
			ret[k] = v
		}
		return ret
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("In", func(t *testing.T) {
				input := copyMap(test.AttrsIn)
				is := &terraform.InstanceState{
					Attributes: input,
				}
				r := &schema.Resource{Schema: test.Schema}
				FixupAsSingleInstanceStateIn(is, r)
				if !cmp.Equal(is.Attributes, test.AttrsOut) {
					t.Errorf("wrong result\ninput: %#v\ngot:   %#v\nwant:  %#v\n\n%s", input, is.Attributes, test.AttrsOut, cmp.Diff(test.AttrsOut, is.Attributes))
				}
			})
			t.Run("Out", func(t *testing.T) {
				input := copyMap(test.AttrsOut)
				is := &terraform.InstanceState{
					Attributes: input,
				}
				r := &schema.Resource{Schema: test.Schema}
				fixupAsSingleInstanceStateOut(is, r)
				if !cmp.Equal(is.Attributes, test.AttrsIn) {
					t.Errorf("wrong result\ninput: %#v\ngot:   %#v\nwant:  %#v\n\n%s", input, is.Attributes, test.AttrsIn, cmp.Diff(test.AttrsIn, is.Attributes))
				}
			})
		})
	}
}
