# Tequila: Compare DDD Model with code

## How to use?
### Prepare
* install golang
* set env variable: GOPATH
* install graphviz(http://graphviz.org/)
* install doxygen(http://www.stack.nl/~dimitri/doxygen/)

### Example:
* DDD Model(![dot file](/examples/step2-problem.dot))

![](https://rawgit.com/newlee/tequila/master/examples/step2-problem.png)

* ![Code](/examples/step2-code/code.h)
* ![Doxygen File](/examples/step2-code/Doxyfile)

### Build & Run Cpp example:
* generate doxygen dot files:
    `doxygen examples/step2-code/Doxyfile`
* build & run example:
    `go build && ./tequila `

### Build & Run Java example:
* generate doxygen dot files:
    `doxygen examples/step2-Java/Doxyfile`
* build & run example:
    `go build && ./tequila `
