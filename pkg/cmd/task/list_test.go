// Copyright © 2019 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/tektoncd/cli/pkg/test"
	cb "github.com/tektoncd/cli/pkg/test/builder"
	testDynamic "github.com/tektoncd/cli/pkg/test/dynamic"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
	pipelinetest "github.com/tektoncd/pipeline/test/v1alpha1"
	"gotest.tools/v3/golden"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTaskList_Invalid_Namespace(t *testing.T) {
	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		},
	}
	cs, _ := test.SeedTestData(t, pipelinetest.Data{Namespaces: ns})
	p := &test.Params{Tekton: cs.Pipeline, Kube: cs.Kube}

	task := Command(p)
	output, _ := test.ExecuteCommand(task, "list", "-n", "foo")
	test.AssertOutput(t, "Error: namespaces \"foo\" not found\n", output)
}

func TestTaskList_Empty(t *testing.T) {
	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
	}
	cs, _ := test.SeedTestData(t, pipelinetest.Data{Namespaces: ns})
	cs.Pipeline.Resources = cb.APIResourceList("v1alpha1", []string{"task"})
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client()
	if err != nil {
		fmt.Println(err)
	}
	p := &test.Params{Tekton: cs.Pipeline, Kube: cs.Kube, Dynamic: dynamic}
	task := Command(p)
	output, err := test.ExecuteCommand(task, "list", "-n", "foo")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	test.AssertOutput(t, "No Tasks found\n", output)
}

func TestTaskList_Only_Tasks_v1alpha1(t *testing.T) {
	clock := clockwork.NewFakeClock()

	tasks := []*v1alpha1.Task{
		tb.Task("tomatoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananas", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("apples", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("potatoes", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("onions", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
	}

	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace",
			},
		},
	}

	version := "v1alpha1"
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client(
		cb.UnstructuredT(tasks[0], version),
		cb.UnstructuredT(tasks[1], version),
		cb.UnstructuredT(tasks[2], version),
		cb.UnstructuredT(tasks[3], version),
		cb.UnstructuredT(tasks[4], version),
		cb.UnstructuredT(tasks[5], version),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	cs, _ := test.SeedTestData(t, pipelinetest.Data{Tasks: tasks, Namespaces: ns})
	p := &test.Params{Tekton: cs.Pipeline, Clock: clock, Kube: cs.Kube, Dynamic: dynamic}
	cs.Pipeline.Resources = cb.APIResourceList("v1alpha1", []string{"task"})
	task := Command(p)

	output, err := test.ExecuteCommand(task, "list", "-n", "namespace")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	golden.Assert(t, output, fmt.Sprintf("%s.golden", t.Name()))
}

func TestTaskList_Only_Tasks_v1beta1(t *testing.T) {
	clock := clockwork.NewFakeClock()

	tasks := []*v1alpha1.Task{
		tb.Task("tomatoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananas", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("apples", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("potatoes", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("onionss", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
	}

	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace",
			},
		},
	}

	version := "v1beta1"
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client(
		cb.UnstructuredT(tasks[0], version),
		cb.UnstructuredT(tasks[1], version),
		cb.UnstructuredT(tasks[2], version),
		cb.UnstructuredT(tasks[3], version),
		cb.UnstructuredT(tasks[4], version),
		cb.UnstructuredT(tasks[5], version),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	cs, _ := test.SeedTestData(t, pipelinetest.Data{Tasks: tasks, Namespaces: ns})
	p := &test.Params{Tekton: cs.Pipeline, Clock: clock, Kube: cs.Kube, Dynamic: dynamic}
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})
	task := Command(p)

	output, err := test.ExecuteCommand(task, "list", "-n", "namespace")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	golden.Assert(t, output, fmt.Sprintf("%s.golden", t.Name()))
}

func TestTaskList_Only_Tasks_no_headers_v1beta1(t *testing.T) {
	clock := clockwork.NewFakeClock()

	tasks := []*v1alpha1.Task{
		tb.Task("tomatoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananas", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("apples", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("potatoes", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("onions", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
	}

	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace",
			},
		},
	}

	version := "v1beta1"
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client(
		cb.UnstructuredT(tasks[0], version),
		cb.UnstructuredT(tasks[1], version),
		cb.UnstructuredT(tasks[2], version),
		cb.UnstructuredT(tasks[3], version),
		cb.UnstructuredT(tasks[4], version),
		cb.UnstructuredT(tasks[5], version),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	cs, _ := test.SeedTestData(t, pipelinetest.Data{Tasks: tasks, Namespaces: ns})
	p := &test.Params{Tekton: cs.Pipeline, Clock: clock, Kube: cs.Kube, Dynamic: dynamic}
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})
	task := Command(p)

	output, err := test.ExecuteCommand(task, "list", "-n", "namespace", "--no-headers")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	golden.Assert(t, output, fmt.Sprintf("%s.golden", t.Name()))
}

func TestTaskList_Only_Tasks_all_namespaces_v1beta1(t *testing.T) {
	clock := clockwork.NewFakeClock()

	tasks := []*v1alpha1.Task{
		tb.Task("tomatoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananas", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("apples", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("potatoes", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("onions", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
		tb.Task("tomates", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangues", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananes", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("pommes", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("patates", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("oignons", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
	}

	version := "v1beta1"
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client(
		cb.UnstructuredT(tasks[0], version),
		cb.UnstructuredT(tasks[1], version),
		cb.UnstructuredT(tasks[2], version),
		cb.UnstructuredT(tasks[3], version),
		cb.UnstructuredT(tasks[4], version),
		cb.UnstructuredT(tasks[5], version),
		cb.UnstructuredT(tasks[6], version),
		cb.UnstructuredT(tasks[7], version),
		cb.UnstructuredT(tasks[8], version),
		cb.UnstructuredT(tasks[9], version),
		cb.UnstructuredT(tasks[10], version),
		cb.UnstructuredT(tasks[11], version),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	cs, _ := test.SeedTestData(t, pipelinetest.Data{Tasks: tasks})
	p := &test.Params{Tekton: cs.Pipeline, Clock: clock, Kube: cs.Kube, Dynamic: dynamic}
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})
	task := Command(p)

	output, err := test.ExecuteCommand(task, "list", "--all-namespaces")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	golden.Assert(t, output, fmt.Sprintf("%s.golden", t.Name()))
}

func TestTaskList_Only_Tasks_all_namespaces_no_headers_v1beta1(t *testing.T) {
	clock := clockwork.NewFakeClock()

	tasks := []*v1alpha1.Task{
		tb.Task("tomatoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangoes", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananas", tb.TaskNamespace("namespace"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("apples", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("potatoes", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("onions", tb.TaskNamespace("namespace"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
		tb.Task("tomates", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-1*time.Minute))),
		tb.Task("mangues", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-20*time.Second))),
		tb.Task("bananes", tb.TaskNamespace("espace-de-nom"), cb.TaskCreationTime(clock.Now().Add(-512*time.Hour))),
		tb.Task("pommes", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("")), cb.TaskCreationTime(clock.Now().Add(-513*time.Hour))),
		tb.Task("patates", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("a test task")), cb.TaskCreationTime(clock.Now().Add(-514*time.Hour))),
		tb.Task("oignons", tb.TaskNamespace("espace-de-nom"), tb.TaskSpec(tb.TaskDescription("a test task to test description of task")), cb.TaskCreationTime(clock.Now().Add(-515*time.Hour))),
	}

	version := "v1beta1"
	tdc := testDynamic.Options{}
	dynamic, err := tdc.Client(
		cb.UnstructuredT(tasks[0], version),
		cb.UnstructuredT(tasks[1], version),
		cb.UnstructuredT(tasks[2], version),
		cb.UnstructuredT(tasks[3], version),
		cb.UnstructuredT(tasks[4], version),
		cb.UnstructuredT(tasks[5], version),
		cb.UnstructuredT(tasks[6], version),
		cb.UnstructuredT(tasks[7], version),
		cb.UnstructuredT(tasks[8], version),
		cb.UnstructuredT(tasks[9], version),
		cb.UnstructuredT(tasks[10], version),
		cb.UnstructuredT(tasks[11], version),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	cs, _ := test.SeedTestData(t, pipelinetest.Data{Tasks: tasks})
	p := &test.Params{Tekton: cs.Pipeline, Clock: clock, Kube: cs.Kube, Dynamic: dynamic}
	cs.Pipeline.Resources = cb.APIResourceList(version, []string{"task"})
	task := Command(p)

	output, err := test.ExecuteCommand(task, "list", "--all-namespaces", "--no-headers")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	golden.Assert(t, output, fmt.Sprintf("%s.golden", t.Name()))
}
