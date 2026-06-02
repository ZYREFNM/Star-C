# Star-Clang
Star-C is a transpiled language written in go that's transpile to C
It uses a GCC implementation to run but that may changes to run with the Clang compiler

- ## Version
	This project is actually instable and doesn't have realeses yet!
	Currently is in Pre-Alpha version 1.0.0

- ## Exemple
	Simple Hello, world!
	```go
    package main;
    import "IO";
    
    func int main() {
        IO::Output.printf("%s", "Hello, world");
        return 0;
    }
    ```
