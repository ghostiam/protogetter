# Protogolint
Welcome to the Protogolint project!

## Overview
Protogolint is a linter developed specifically for Go programmers working with nested `protobuf` types.\
It's designed to aid developers in preventing `invalid memory address or nil pointer dereference` errors arising from direct access of nested `protobuf` fields.

When working with `protobuf`, it's quite common to have complex structures where a message field is contained within another message, which itself can be part of another message, and so on.
If these fields are accessed directly and some field in the call chain will not be initialized, it can result in application panic.

Protogolint addresses this issue by suggesting use of getter methods for field access.

## How does it work?
Protogolint analyzes your Go code and helps detect direct `protobuf` field accesses that could give rise to panic.\
The linter suggests using getters:
```go
m.GetFoo().GetBar().GetBaz()
```
instead of direct field access:
```go
m.Foo.Bar.Baz
```

And you will then only need to perform a nil check after the final call:
```go
if m.GetFoo().GetBar().GetBaz() != nil {
    // Do something with m.GetFoo().GetBar().GetBaz()
}
```
instead of:
```go
if m.Foo != nil {
    if m.Foo.Bar != nil {
        if m.Foo.Bar.Baz != nil {
            // Do something with m.Foo.Bar.Baz
        }
    }
}
```
which simplifies the code and makes it more reliable.

## Installation

```bash
go install github.com/ghostiam/protogolint@latest
```

## Usage

To run the linter:
```bash
protogolint ./...
```

Or to apply suggested fixes directly:
```bash
protogolint --fix ./...
```
