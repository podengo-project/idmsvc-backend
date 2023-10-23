# Design notes

It was tried to use the following guidelines:

- Single responsibility principla: Each component has one, and only one reason to exist.
  - Reduce the complexity for one component defining a small scope.
  - Make the unit tests easy to implement as the scope of the component is reduced.
- Interface segregation: Define behavior with interfaces, and try to define them
  simple, so they can be combined if necessary to define more comples behaviors.
  - This is important for the unit tests when using mocks; mockery tool can generate the mocks for all the defined interfaces, which reduce a lot of boilerplate.
- Liskov Substitution: Components that implements the same interface can be exchangable. This has been specilly in mind when defining the Service interface to start/stop different macro-components.
- Dependency inversion: Depends on interfaces instead of concret implementations. To achieve this, interfaces and implementation should be implemented in different packages. If an interface has several implementations, make sense to 
