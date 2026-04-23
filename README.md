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
    
    Bank exemple (this exemple does not intend to be an accurate Bank system program but instead a functionalities show up)
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
            this.money = money;
        }
        
        func void activate() {
            this.accessible = true;
        }
        
        func void deactivate() {
            this.accessible = false;
        }
        
        func int main() {
            var Bank newBank = BankAccount.new();
            newBank.activate();
            if (newBank.money <= 0) {
                newBank.deactivate();
            }
            return 0;
        }
    }