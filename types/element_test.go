package types_test

import (
	"encoding/json"
	"testing"

	"github.com/ideatru/crosschain/types"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToElements(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    string
		Expected types.Elements
	}

	testcases := []Testcase{
		{
			Name: "single-type",
			Input: `[
				{
					"type": "string"
				}
			]`,
			Expected: types.Elements{
				{
					Type: types.String,
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
			Expected: types.Elements{
				{
					Type: types.Number,
				},
				{
					Type: types.Array,
					Children: types.Elements{
						{
							Type: types.String,
						},
					},
				},
				{
					Type: types.Object,
					Children: types.Elements{
						{
							Name: "field-string",
							Type: types.String,
						},
						{
							Name: "field-number",
							Type: types.Number,
						},
						{
							Name: "field-boolean",
							Type: types.Bool,
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			var elements types.Elements

			err := json.Unmarshal([]byte(tc.Input), &elements)
			assert.NoError(t, err)

			assert.Equal(t, tc.Expected, elements)
		})
	}
}
