// Package builder to help on building data models.
//
// This package implement several builders to reduce boileplate
// when we are coding tests for the service, which include any
// kind of tests.
//
// TL;DR
//
// Builder pattern is used.
//
// Long history, an initial builder pattern with a wrong approach
// was considered initially, where for each field of the wrapped
// struct to build, it was added a flag field to don't generate
// information for that fields; it was discarded because it was
// adding too much boilerplate to implement new builders.
//
// After was considered a different implementation by using the
// golang options patter, which was generating some random data
// for the wrapped builder struct, and overriding the information
// by calling the slice of options funtions provided.
// This was discarded because still the golang options pattern
// was adding too much boilerplate to add new builders; additionally
// it was evoking name conflicts with the options for other new
// strutcs, so that was makin a bit difficult to get good semantic
// as the builder functionality is increased for all the different
// structs.
//
// Finally it was get back to the Build pattern, but removing the
// flags for each provided field which reduce the boilerplate; with
// this approach we create a random initial object, override the
// fields to customize. It is easy to compose complex builder by
// using composition, and the semantic allow to use the builder
// in a very readable way, letting the code to speak by itself.
//
// ```golang
//
//	type MyType interface {
//		Build() model.MyType
//		WithFieldName1(value int64) MyType
//	}
//
//	type myType struct {
//		model.MyType
//	}
//
//	func NewMyType() MyType {
//		return &model.MyType {
//			FieldName1: GenRandNum(0, 10000),
//		}
//	}
//
//	func (b *myType) Buid() MyTpye {
//		return b.MyType
//	}
//
//	func (b *myType) WithFieldName1(value int64) MyType {
//		b.MyType.FiledName1 = value
//		return b
//	}
//
// ```
//
// The below would be the minimal; we can see that generate some
// random data, and we override the fields as we call method
// members of the builder, finally we call Build() method to get
// the generated data. Now the above example would be used as the
// below.
//
// ```golang
//
//	func TestMyTypeCase1(t *testify.Testing) {
//		o := builder.NewMyType().WithFieldName1(45)
//	}
//
// ```
//
// Into the above leverage example we see how the semantic
// referencing the package and using the specific factory
// make the code speak by itself.
//
// TODO The current approach match a repetitive pattern for every
// builder that maybe a tool could implement to generate the code
// given the gorm database model; but this would be for the future.
// The tool indicated would allow to generate the builder boilerplate
// and keep it on sync with the model changes.
package builder
