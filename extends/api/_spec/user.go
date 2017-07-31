package _spec

import (
	. "github.com/mbict/gogen/dsl"
	. "github.com/mbict/gogen/extends/api/dsl"
)

var Username = Type("Username", String, func() {
	Description("Email/ Username of the User")
	//Format(Email)
})

var User = Type("User", func() {
	Namespace("user")
	Description("User")

	Attribute("user_id", UUID)
	Attribute("username", Username)
	Attribute("password", String)
	Attribute("active", Boolean, "Active and enabled user")
	Attribute("last_login_at", DateTime, "Last time the user logged in")
	Attribute("created_at", DateTime, "Creation date")
	Attribute("updated_at", DateTime, "Last date when the user details where changed")

	Required("user_id", "username", "password", "active", "last_login_at", "created_at", "updated_at")
})

var UserService = Service("User", func() {
	Namespace("user")

	HTTP(func() {
		Path("/users")
	})

	Error("not_found", String, "No user found matching your criteria")
	Error("duplicate_username", String, "The username is already in use")
	Error("invalid_username", String, "The username is invalid, use a email address")

	Method("GetUser", func() {
		Description("Retrieve a user")
		Result(User)
		Payload(func() {
			Attribute("user_id", UUID)
		})

		HTTP(func() {
			GET("/{user_id}")
			//Response("not_found", NotFound)
		})
	})

	Method("GetUsers", func() {
		Description("List all users")
		//Result(func() {
		//	Attribute("users", ArrayOf(User), "Collection of users")
		//})
		Result(ArrayOf(User))

		HTTP(func() {
			GET("")
		})
	})

	Method("Authenticate", func() {
		Description("Authenticate a username and password combination")
		Error("unauthorized", String, "Cannot authenticate the username and password combination")
		Result(Boolean)
		Payload(func() {
			Required("username", "password")
			Attribute("username", Username)
			Attribute("password", String)
		})
	})

	Method("CreateUser", func() {
		Description("Create a new user")
		Payload(func() {
			Required("username", "password")
			Attribute("username", Username)
			Attribute("password", String, "Plain password")
			Attribute("active", Boolean, "Is the user active/enabled")
		})
		HTTP(func() {
			POST("")
			//Response("invalid_username", BadRequest)
			//Response("duplicate_username", BadRequest)
			//Response(StatusCreated)
		})
	})

	Method("UpdateUser", func() {
		Description("Update user")
		Payload(func() {
			Required("user_id")
			Attribute("user_id", UUID)
			Attribute("username", Username)
			Attribute("password", String)
			Attribute("active", Boolean, "Is the user active/enabled")
		})

		HTTP(func() {
			PUT("/{user_id}")
			//Response("not_found", NotFound)
			//Response("invalid_username", BadRequest)
			//Response("duplicate_username", BadRequest)
			//Response(StatusNoContent)
		})
	})

	Method("DeleteUser", func() {
		Description("Delete a user")
		Payload(func() {
			Attribute("user_id", UUID)
		})
		HTTP(func() {
			DELETE("/{user_id}")
			//Error("not_found", NotFound)
			//Response(StatusNoContent)
		})
	})
})

