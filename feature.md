
## Core Element
### Package
* structure slice (ref)

### structure
* full name (package name + self name) as identify
* property (native type)
* other structure reference (full name / identity)
* method slice (full name / identity)

### method
* full name / identity (??)
* parameters slice (aggregate)
    * name
    * type
* structure dependency
    * full name
    * method
    * property (contains other structure reference), maybe have cascade like a.b.c.d or a->b->c->d
* return value
    * name and type (some language support more than one return value)    