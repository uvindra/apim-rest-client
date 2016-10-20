package swagger



func getResponseCode() *Code_200 {
	code := Code_200{}
	code.Code_200 = &Description{}

	return &code
}

func getMethodProperties(response *Code_200, description string, summary string) *MethodProperties  {
	properties := MethodProperties{}

	properties.Responses = response
	properties.ResourceDescription = description
	properties.Summary = summary
	properties.Consumes = []string{"application/json"}
	properties.Produces = []string{"application/json"}

	return &properties
}

func getPaths() *PathsType {
	read := ReadResource{}
	read.Read = getMethodProperties(getResponseCode(), "Read Data", "Read")

	create := CreateResource{}
	create.Create = getMethodProperties(getResponseCode(), "Create Data", "Create")

	remove := RemoveResource{}
	remove.Remove = getMethodProperties(getResponseCode(), "Remove Data", "Remove")

	replace := ReplaceResource{}
	replace.Replace = getMethodProperties(getResponseCode(), "Replace Data", "Replace")

	update := UpdateResource{}
	update.Update = getMethodProperties(getResponseCode(), "Update Data", "Update")

	check := CheckResource{}
	check.Check = getMethodProperties(getResponseCode(), "Check if resource exists", "check")

	what := WhatResource{}
	what.What = getMethodProperties(getResponseCode(), "What options are available?", "what")

	paths := PathsType{}
	paths.Read = &read
	paths.Create = &create
	paths.Remove = &remove
	paths.Replace = &replace
	paths.Update = &update
	paths.Check = &check
	paths.What = &what

	return &paths
}

func GetSwagger(name string, version string) *Swagger {
	info := InfoType{}
	info.Title = name
	info.Version = version

	swagger := Swagger{}
	swagger.Version = "2.0"
	swagger.Info = &info
	swagger.Paths = getPaths()

	return &swagger
}


