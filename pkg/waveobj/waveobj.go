// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

// Package waveobj defines the core object model for WaveTerm.
package waveobj

import (
	"fmt"
	"reflect"
)

// OType represents the object type identifier.
type OType = string

// OID represents a unique object identifier.
type OID = string

// WaveObj is the base interface all wave objects must implement.
type WaveObj interface {
	GetOType() OType
	GetOID() OID
}

// MetaType holds metadata key-value pairs for wave objects.
type MetaType map[string]any

// WaveObjBase provides a base implementation of WaveObj.
type WaveObjBase struct {
	OType OType  `json:"otype"`
	OID   OID    `json:"oid"`
	Meta  MetaType `json:"meta,omitempty"`
}

func (b *WaveObjBase) GetOType() OType {
	return b.OType
}

func (b *WaveObjBase) GetOID() OID {
	return b.OID
}

// registry maps OType strings to their reflect.Type for deserialization.
var registry = map[OType]reflect.Type{}

// RegisterType registers a WaveObj type with the given OType string.
// T must be a pointer to a struct that embeds WaveObjBase.
func RegisterType[T WaveObj](otype OType) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	registry[otype] = t
}

// MakeObjByType creates a new zero-value WaveObj for the given OType.
func MakeObjByType(otype OType) (WaveObj, error) {
	t, ok := registry[otype]
	if !ok {
		return nil, fmt.Errorf("unknown otype: %q", otype)
	}
	obj := reflect.New(t).Interface()
	wobj, ok := obj.(WaveObj)
	if !ok {
		return nil, fmt.Errorf("registered type %q does not implement WaveObj", otype)
	}
	return wobj, nil
}

// GetRegisteredTypes returns a sorted list of all registered OTypes.
func GetRegisteredTypes() []OType {
	result := make([]OType, 0, len(registry))
	for k := range registry {
		result = append(result, k)
	}
	return result
}
