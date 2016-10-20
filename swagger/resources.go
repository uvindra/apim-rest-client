package swagger

type ReadResource struct {
	Read *MethodProperties `json:"get"`
}


type CreateResource struct {
	Create *MethodProperties `json:"post"`
}

type RemoveResource struct {
	Remove *MethodProperties `json:"delete"`
}


type ReplaceResource struct {
	Replace *MethodProperties `json:"put"`
}

type UpdateResource struct {
	Update *MethodProperties `json:"patch"`
}


type CheckResource struct {
	Check *MethodProperties `json:"head"`
}

type WhatResource struct {
	What *MethodProperties `json:"options"`
}