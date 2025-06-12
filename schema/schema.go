package schema

/*
value: ["lorem, ipsum"]

schema: [
	{
		"type": "string",
	}
]
*/

/*
value: [
	100,
	[
		"value-1",
		"value-2"
	],
	{
		"field-string": "lorem, ipsum",
		"field-number": 99.99,
		"field-boolean": true
	}
]

schema: [
	{
		"type": "number",
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
]
*/

type Element struct {
	Name     string      `json:"name"`
	Ty       ElementType `json:"type"`
	Nullable bool        `json:"nullable"`
	Children []Element   `json:"children,omitempty"`
}
