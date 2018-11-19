/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	devv1 "github.com/proepkes/speeddate/src/pkg/apis/dev/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGameServers implements GameServerInterface
type FakeGameServers struct {
	Fake *FakeDevV1
	ns   string
}

var gameserversResource = schema.GroupVersionResource{Group: "dev.speeddate", Version: "v1", Resource: "gameservers"}

var gameserversKind = schema.GroupVersionKind{Group: "dev.speeddate", Version: "v1", Kind: "GameServer"}

// Get takes name of the gameServer, and returns the corresponding gameServer object, and an error if there is any.
func (c *FakeGameServers) Get(name string, options v1.GetOptions) (result *devv1.GameServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(gameserversResource, c.ns, name), &devv1.GameServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*devv1.GameServer), err
}

// List takes label and field selectors, and returns the list of GameServers that match those selectors.
func (c *FakeGameServers) List(opts v1.ListOptions) (result *devv1.GameServerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(gameserversResource, gameserversKind, c.ns, opts), &devv1.GameServerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &devv1.GameServerList{ListMeta: obj.(*devv1.GameServerList).ListMeta}
	for _, item := range obj.(*devv1.GameServerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gameServers.
func (c *FakeGameServers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(gameserversResource, c.ns, opts))

}

// Create takes the representation of a gameServer and creates it.  Returns the server's representation of the gameServer, and an error, if there is any.
func (c *FakeGameServers) Create(gameServer *devv1.GameServer) (result *devv1.GameServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(gameserversResource, c.ns, gameServer), &devv1.GameServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*devv1.GameServer), err
}

// Update takes the representation of a gameServer and updates it. Returns the server's representation of the gameServer, and an error, if there is any.
func (c *FakeGameServers) Update(gameServer *devv1.GameServer) (result *devv1.GameServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(gameserversResource, c.ns, gameServer), &devv1.GameServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*devv1.GameServer), err
}

// Delete takes name of the gameServer and deletes it. Returns an error if one occurs.
func (c *FakeGameServers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(gameserversResource, c.ns, name), &devv1.GameServer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGameServers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(gameserversResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &devv1.GameServerList{})
	return err
}

// Patch applies the patch and returns the patched gameServer.
func (c *FakeGameServers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *devv1.GameServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(gameserversResource, c.ns, name, pt, data, subresources...), &devv1.GameServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*devv1.GameServer), err
}