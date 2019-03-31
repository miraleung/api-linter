package lint

import (
	"context"

	"github.com/golang/protobuf/reflect/protoreflect"
)

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule interface {
	// ID returns the unique `RuleID` for this rule.
	ID() RuleID
	// Description returns a short summary about this rule.
	Description() string
	// URL returns a link at which the full documentation
	// about this rule can be accessed.
	URL() string
	// FileTypes returns the types of files that this rule is targeting to.
	// E.g., `ProtoFile` for protobuf files.
	FileTypes() []FileType
	// Category returns the category of findings produced
	// by this rule, e.g. Problem, Suggestion, etc.
	Category() Category
	// Lint performs the linting process.
	Lint(Request) (Response, error)
}

// RuleID describes a `Rule` ID.
type RuleID struct {
	// Set is where a rule belongs to and usually is the scope
	// of the contained rules should be applied to. For example,
	// the `core` set contains all rules that should be applied to
	// all Google APIs, while the set `google.corpeng` contains rules
	// that should be applied only to Google Corp Eng APIs.
	Set string
	// Name is the unique name in the set.
	Name string
}

// FileType defines a file type that a rule is targeting to.
type FileType string

const (
	// ProtoFile indicates that the targeted file is a protobuf-definition file.
	ProtoFile FileType = "proto-file"
)

// Category defines the category of the findings produced by a rule.
type Category string

const (
	// CategoryError indicates that in the API, something will cause errors.
	CategoryError Category = "API-Linter-Error"
	// CategorySuggestion indicates that in the API, something can be do better.
	CategorySuggestion Category = "API-Linter-Suggestion"
)

// Context provides additional information for a rule to perform linting.
type Context struct {
	context    context.Context
	descSource DescriptorSource
}

// NewContext creates a new `Context`.
func NewContext(ctx context.Context) Context {
	return Context{
		context: ctx,
	}
}

// NewContextWithDescriptorSource creates a new `Context` with the source.
func NewContextWithDescriptorSource(ctx context.Context, source DescriptorSource) Context {
	return Context{
		context:    ctx,
		descSource: source,
	}
}

// DescriptorSource returns a `DescriptorSource` if available; otherwise,
// returns (nil, ErrSourceInfoNotAvailable).
//
// The returned `DescriptorSource` contains additional information
// about a protobuf descriptor, such as comments and location lookups.
func (c Context) DescriptorSource() DescriptorSource {
	return c.descSource
}

// Request defines input data for a rule to perform linting.
type Request interface {
	// ProtoFile returns a `FileDescriptor` when the rule's `FileTypes`
	// contains `ProtoFile`.
	ProtoFile() protoreflect.FileDescriptor
	Context() Context
}

// Response describes the result returned by a rule.
type Response struct {
	Problems []Problem
}

// Merge merges another response.
func (resp *Response) Merge(other Response) {
	resp.Problems = append(resp.Problems, other.Problems...)
}
