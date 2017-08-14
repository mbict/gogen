package basic

import (
	. "github.com/mbict/gogen/dsl"
)

//Primitive user type defined from other primitive type
var UserId = Type("UserId", String, "Custom user type")

//Structure type
var UserStructure = Type("UserId", String, "Custom user type", func() {

	//Basic attribute
	Attribute("username", String)

	//Attribute with user type
	Attribute("user_id", UserId)

	//Attribute with additional DSL (validation)
	Attribute("Password", String, func() {
		MinLength(8)
		MaxLength(128)
	})

	//Attribute with description
	Attribute("notes", String, "A description for the notes attribute")

	//Tag with description and extra validation defined in the DSL
	Attribute("tag", String, "Description for username", func() {
		Pattern("^[a-zA-Z0-9_.\\-]*$")
		MinLength(3)
		MaxLength(64)
	})
})
