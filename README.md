# Star-Clang
Star-C is a transpiled language written in go that's transpile to C
It uses a GCC implementation to run but that may changes to run with the Clang compiler

- ## Version
	This project is actually instable and doesn't have realeses yet!
	Currently is in Pre-Alpha version 1.2.0

- ## Exemples
	Simple Hello, world!
	```go
    package main;
    import "IO";
    
    func int main() {
        IO::Output.printf("%s", "Hello, world");
        return 0;
    }
    ```
    Change a player's username in a game and print it
    ```go
    package main;
    import "IO";
    
    class Player {
        var <get, set> string name;
    }
    
    func int main() {
        var Player newPlayer;
        newPlayer.set<name("AlphaPlayer111")>;
        IO::Output.printfln("%s", newPlayer.get<name()>);
        return 0;
    }
    ```
    ```bash
    Output: AlphaPlayer111
