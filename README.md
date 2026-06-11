# Star-Clang
Star-C is a transpiled language written in go that's transpile to C
It uses a GCC implementation to run but that may changes to run with the Clang compiler

- ## Philosophy
    Star-C has been created to be a fast and powerful language that aim:
    - ### Lightweight:
        Star-C aims to be as light as possible while not decreasing performances.
    - ### Flexibility:
        Star-C's functionnalities have been made to be able to work for as much case possible in as much context possible.

- ## Version
	This project is actually instable and doesn't have realeses yet!
	Currently is in Pre-Alpha version 1.3.0

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
