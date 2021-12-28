package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

moduleField: schema.#resource & {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	rbac: {
		operations: {
			"recod.value.read": description:   "Read field value on records"
			"recod.value.update": description: "Update field value on records"
		}
	}

	locale: {
		skipSvc: true

		keys: {
			label: {}
			descriptionView: {
				path: ["meta", "description", "view"]
				customHandler: true
			}
			descriptionEdit: {
				path: ["meta", "description", "edit"]
				customHandler: true
			}
			hintView: {
				path: ["meta", "hint", "view"]
				customHandler: true
			}
			hintEdit: {
				path: ["meta", "hint", "edit"]
				customHandler: true
			}
			validatorError: {
				path: ["expression", "validator", {part: "validatorID", var: true}, "error"]
				customHandler: true
			}
			optionsOptionTexts: {
				path: ["meta", "options", {part: "value", var: true}, "text"]
				customHandler: true
			}
		}
	}
}
