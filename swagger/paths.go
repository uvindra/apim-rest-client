package swagger

type PathsType struct {
	Read *ReadResource  `json:"/read"`
	Create *CreateResource `json:"/create"`
	Remove *RemoveResource  `json:"/remove"`
	Replace *ReplaceResource `json:"/replace"`
	Update *UpdateResource    `json:"/update"`
	Check *CheckResource      `json:"/check"`
	What *WhatResource        `json:"/what"`
}
