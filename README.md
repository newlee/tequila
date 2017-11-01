# Tequila: Compare DDD Model with code

## How to use?
### Prepare
* install golang
* set env variable: GOPATH
* install graphviz(http://graphviz.org/)
* install doxygen(http://www.stack.nl/~dimitri/doxygen/)

### Example:
* DDD Model(![dot file](/examples/cargo-problem.dot))

![](https://rawgit.com/newlee/tequila/master/examples/cargo-problem.png)

* ![Code](/examples/bc-code)
* ![Doxygen File](/examples/bc-code/Doxyfile)

### Features Done
* DDD model validate
* Cpp code check with DDD model
* Inherit support
* Include dependencies visualization

### Features TODO
* Output detail result for cpp code check
* Cpp code check with DDD solution domain

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
