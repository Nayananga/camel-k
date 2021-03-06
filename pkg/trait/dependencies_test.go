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

package trait

import (
	"testing"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util/test"

	"github.com/stretchr/testify/assert"
)

func TestDependenciesTraitApplicability(t *testing.T) {
	e := &Environment{
		Integration: &v1alpha1.Integration{},
	}

	trait := newDependenciesTrait()
	enabled, err := trait.Configure(e)
	assert.Nil(t, err)
	assert.False(t, enabled)

	e.Integration.Status.Phase = v1alpha1.IntegrationPhaseNone
	enabled, err = trait.Configure(e)
	assert.Nil(t, err)
	assert.False(t, enabled)

	e.Integration.Status.Phase = v1alpha1.IntegrationPhaseInitialization
	enabled, err = trait.Configure(e)
	assert.Nil(t, err)
	assert.True(t, enabled)
}

func TestIntegrationDefaultDeps(t *testing.T) {
	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	e := &Environment{
		CamelCatalog: catalog,
		Integration: &v1alpha1.Integration{
			Spec: v1alpha1.IntegrationSpec{
				Sources: []v1alpha1.SourceSpec{
					{
						DataSpec: v1alpha1.DataSpec{
							Name:    "Request.java",
							Content: `from("direct:foo").to("log:bar");`,
						},
						Language: v1alpha1.LanguageJavaSource,
					},
				},
			},
			Status: v1alpha1.IntegrationStatus{
				Phase: v1alpha1.IntegrationPhaseInitialization,
			},
		},
	}

	trait := newDependenciesTrait()
	enabled, err := trait.Configure(e)
	assert.Nil(t, err)
	assert.True(t, enabled)

	err = trait.Apply(e)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"camel:core", "camel:direct", "camel:log", "runtime:jvm"}, e.Integration.Status.Dependencies)
}

func TestIntegrationCustomDeps(t *testing.T) {
	catalog, err := test.DefaultCatalog()
	assert.Nil(t, err)

	e := &Environment{
		CamelCatalog: catalog,
		Integration: &v1alpha1.Integration{
			Spec: v1alpha1.IntegrationSpec{
				Dependencies: []string{
					"camel:undertow",
					"org.foo:bar",
				},
				Sources: []v1alpha1.SourceSpec{
					{
						DataSpec: v1alpha1.DataSpec{
							Name:    "Request.java",
							Content: `from("direct:foo").to("log:bar");`,
						},
						Language: v1alpha1.LanguageJavaSource,
					},
				},
			},
			Status: v1alpha1.IntegrationStatus{
				Phase: v1alpha1.IntegrationPhaseInitialization,
			},
		},
	}

	trait := newDependenciesTrait()
	enabled, err := trait.Configure(e)
	assert.Nil(t, err)
	assert.True(t, enabled)

	err = trait.Apply(e)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"camel:core", "camel:direct", "camel:log",
		"camel:undertow", "org.foo:bar", "runtime:jvm"}, e.Integration.Status.Dependencies)
}
