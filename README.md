:# Star-Clang
Star-C is a source-to-source compiled language written in go that uses a C backend.
It uses a GCC implementation to run but that may changes to run with the Clang compiler (or both).

- ## Philosophy
    - ## Niche
        Star-C is a system language, in the same niche as C, C++, Rust and Zig.

- ## Integrations:
    - ### Memory:
        Star-C manages memory using an Arena but also allows manual control. Star-C also have the SuperNova system, a precompiler letting the user interact with memory AOT.

- ## Version
	This project is actually instable and doesn't have realeses yet!
	Currently is in Pre-Alpha version 2.1.1

- ## Exemples
	Simple Hello, world!
	```go
    module helloworld;
    import "IO";
    
    func int main() {
        IO::printf("%s", "Hello, world");
        return 0;
    }
    ```
    Change a player's username in a game and print it
    ```go
    module username;
    import "IO";
    
    typedef struct Player {
        var (get, set) string name;
    }
    
    func int main() {
        var Player newPlayer;
        set { newPlayer.name("AlphaPlayer111") }
        IO::printfln("%s", newPlayer.get{name()});
        return 0;
    }
    ```
    ```bash
    Output: AlphaPlayer111