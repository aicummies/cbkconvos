# cbkconvos
get your conversations in a tavern compatible way :)

## usage

1. install go
2. put your api token in the `bearer` file first line
3. run the command this way `go run . -c <yourConversationId> --charactername <yourCharacterName>`

example:
```
$ go run . -c c1ff34u3847382rd --charactername Megumin
```
This script creates two files:
- `backstory_167984154643401.txt` contains the gaslight + character info that you inputted in
- `167984154643401.jsonl` contains the actual messages in a format compatible with tavern
