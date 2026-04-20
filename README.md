# Star-C
Star-C is a modern C like language that tries to offer and prioritize a minimalistic and lightweight tool but with fast and powerful features

- ## Version
	Currently is Pre-Alpha version 0.5.4

- ## Exemple
	Simple Hello, world!
	```go
    func int main() {
        print "Hello, world!";
        return 0;
    }
    ```
    
    Advanced bank account representation
    ``` go
    class BankAccount {
        public var string accountName;
        private {
            var boolean accessible = true;
            var <get, set> float64 money {
                get() <- this;
                set(float64 value) <- this;
            }
        }
        
        func BankAccount new(float money) {
            this.accessible = true;
            this.money = money;
        }
        
        func void deactivate() {
            this.accessible = false;
        }
    }