/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package build

// Request represent a request to build a specific code
type Request struct {
	Identifier   Identifier
	Code         Source
	Dependencies []string
}

// Identifier identifies a build
type Identifier struct {
	Name      string
	Qualifier string
}

// Source represent the integration code
type Source struct {
	Name     string
	Content  string
	Language string
}

// Result represents the result of a build
type Result struct {
	Source *Request
	Status Status
	Image  string
	Error  error
}

// Builder is supertype of all builders
type Builder interface {
	Build(Request) <-chan Result
}

// Status --
type Status int

const (
	// StatusNotRequested --
	StatusNotRequested Status = iota

	// StatusStarted --
	StatusStarted

	// StatusCompleted --
	StatusCompleted

	// StatusError --
	StatusError
)