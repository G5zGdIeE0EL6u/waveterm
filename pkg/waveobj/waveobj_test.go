// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

package waveobj_test

import (
	"testing"

	"github.com/wavetermdev/waveterm/pkg/waveobj"
)

// testObj is a simple WaveObj implementation for testing.
type testObj struct {
	OType string `json:"otype"`
	OID   string `json:"oid"`
	Name  string `json:"name"`
}

func (t *testObj) GetOType() string    { return "test" }
func (t *testObj) GetOID() string      { return t.OID }
func (t *testObj) SetOID(oid string)   { t.OID = oid }

// anotherObj is a second WaveObj implementation for testing multiple registrations.
type anotherObj struct {
	OType string `json:"otype"`
	OID   string `json:"oid"`
	Value int    `json:"value"`
}

func (a *anotherObj) GetOType() string  { return "another" }
func (a *anotherObj) GetOID() string    { return a.OID }
func (a *anotherObj) SetOID(oid string) { a.OID = oid }

func TestRegisterType(t *testing.T) {
	waveobj.RegisterType[*testObj]()
	types := waveobj.GetRegisteredTypes()
	found := false
	for _, tp := range types {
		if tp == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected type 'test' to be registered, got: %v", types)
	}
}

func TestRegisterMultipleTypes(t *testing.T) {
	waveobj.RegisterType[*testObj]()
	waveobj.RegisterType[*anotherObj]()
	types := waveobj.GetRegisteredTypes()
	expected := map[string]bool{"test": false, "another": false}
	for _, tp := range types {
		if _, ok := expected[tp]; ok {
			expected[tp] = true
		}
	}
	for tp, found := range expected {
		if !found {
			t.Errorf("expected type '%s' to be registered", tp)
		}
	}
}

func TestMakeObjByType(t *testing.T) {
	waveobj.RegisterType[*testObj]()
	obj, err := waveobj.MakeObjByType("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatal("expected non-nil object")
	}
	if obj.GetOType() != "test" {
		t.Errorf("expected otype 'test', got '%s'", obj.GetOType())
	}
}

func TestMakeObjByTypeUnknown(t *testing.T) {
	_, err := waveobj.MakeObjByType("nonexistent")
	if err == nil {
		t.Error("expected error for unknown type, got nil")
	}
}

// TestMakeObjByTypeAndSetOID verifies that a newly created object can have its OID set correctly.
// Added this to ensure SetOID works as expected after MakeObjByType.
func TestMakeObjByTypeAndSetOID(t *testing.T) {
	waveobj.RegisterType[*testObj]()
	obj, err := waveobj.MakeObjByType("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const testOID = "test-oid-1234"
	obj.SetOID(testOID)
	if obj.GetOID() != testOID {
		t.Errorf("expected OID '%s', got '%s'", testOID, obj.GetOID())
	}
}
