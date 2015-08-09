Debugging Tools for Go
======================
This project implements tools that are helpful when debugging Go
programs running in development, testing, or production environments.

Initially, this project is focused on implementing debugging support for
Go programs running on Microsoft Windows. AS the project develops, I
will be looking at adding more support for Linux and Mac programs too.

Get the Code
------------
Getting the code is easy using the `go get` command. In a terminal or
command prompt window, execute the following command:

    $ go get github.com/mfcollins3/go-debugger

To use the debugging library in your own programs or libraries, add The
following `import` statement at the top of your programs:

```go
import debugger "github.com/mfcollins3/go-debugger"
```

Features
--------
### OutputDebugString Support (Microsoft Windows)

The `go-debugger` library supports writing messages at runtime to an
attached debugger or [DebugView](https://technet.microsoft.com/en-us/Library/bb896647.aspx)
using the Microsoft Windows [OutputDebugString](https://msdn.microsoft.com/en-us/library/windows/desktop/aa363362(v=vs.85).aspx)
API. This feature has been designed to support the [io.Writer](http://golang.org/pkg/io/#Writer)
interface because there is so much support within the Go framework and
third-party libaries for writing to `io.Writer` objects.

Writing to the debugger is easy:

```go
import debugger "github.com/mfcollins3/go-debugger"

func main() {
  debugger.Console.WriteString("The program is starting")

  // TODO: implement the program

  debugger.Console.WriteString("The program is ending")  
}
```

In addition, `debugger.Console` can be used with functions that rely on
`io.Writer` such as the `fmt` package:

```go
import debugger "github.com/mfcollins3/go-debugger"

func main() {
  fmt.Fprint(debugger.Console, "The program is starting")

  var firstName = "Michael", lastName = "Collins"
  fmt.Fprintf(debugger.Console, "Hello, %s %s", firstName, lastName)

  fmt.Fprint(debugger.Console, "The program is ending")
}
```

On unsupported platforms (currently anything that isn't Microsoft
Windows), `debugger.Console` is backed by a null implementation. This
code will run successfully with minimal performance impact on those
platforms.
