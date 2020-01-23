# Entity Component System in Go

## Introduction

This is my attempt to build a high-performance ECS (Entity, Component, System) engine in pure Go. Having investigated other ECS systems which have been written to date, I realised that none of them is really a pure ECS as the data layout is often using pointers or interfaces, which would discard most of the benefits of such a system. 


## Entities

Entities in this system are represented with the `Entity` struct as opposed to a simple identifier (which they also have). The main reason behind this is to support cleanly deleting an entity from all of the component managers.

## Component Generation

The components are protected through a generated "component manager" which can be found in `component` package. In order to generate one, I used [genny](https://github.com/cheekybits/genny) which can be used with go generate as well (see `component/generic.go`).

```
//go:generate genny -in=$GOFILE -out=z_components.go gen "TType="
```

