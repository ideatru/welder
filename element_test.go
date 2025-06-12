package bridge

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToElements(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    string
		Expected []Element
	}

	testcases := []Testcase{
		{
			Name: "single-type",
			Input: `[
				{
					"type": "string"
				}
			]`,
			Expected: []Element{
				{
					Ty: String,
				},
			},
		},
		{
			Name: "multiple-types",
			Input: `[
				{
					"type": "number"
				},
				{
					"type": "array",
					"children": [
						{
							"type": "string"
						}
					]
				},
				{
					"type": "object",
					"children": [
						{
							"name": "field-string",
							"type": "string"
						},
						{
							"name": "field-number",
							"type": "number"
						},
						{
							"name": "field-boolean",
							"type": "boolean"
						}
					]
				}
			]`,
			Expected: []Element{
				{
					Ty: Number,
				},
				{
					Ty: Array,
					Children: []Element{
						{
							Ty: String,
						},
					},
				},
				{
					Ty: Object,
					Children: []Element{
						{
							Name: "field-string",
							Ty:   String,
						},
						{
							Name: "field-number",
							Ty:   Number,
						},
						{
							Name: "field-boolean",
							Ty:   Boolean,
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			var elements []Element
			err := json.Unmarshal([]byte(tc.Input), &elements)
			if err != nil {
				t.Fatalf("failed to unmarshal input: %v", err)
			}

			assert.Equal(t, tc.Expected, elements)
		})
	}
}
