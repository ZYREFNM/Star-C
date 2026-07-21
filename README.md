# Star-Clang
Star-C is a source-to-source compiled language written in go that uses a C backend.
It uses a GCC implementation to run but that may changes to run with the Clang compiler (or both).

- ## Philosophy
    - ## Fill later or smt

- ## Integrations:
    - ### Memory:
        Star-C manages memory using an Arena but also allows manual control.

- ## Version
	This project is actually instable and doesn't have realeses yet!
	Currently is in Pre-Alpha version 1.3.0

- ## Exemples
	Simple Hello, world!
	```go
    module main;
    import "IO";
    
    func int main() {
        IO::out.printf("%s", "Hello, world");
        return 0;
    }
    ```
    Change a player's username in a game and print it
    ```go
    module main;
    import "IO";
    
    class Player {
        var (get, set) string name;
    }
    
    func int main() {
        var Player newPlayer;
        newPlayer.set{name("AlphaPlayer111")};
        IO::out.printfln("%s", newPlayer.get{name()});
        return 0;
    }
    ```
    ```bash
    Output: AlphaPlayer111