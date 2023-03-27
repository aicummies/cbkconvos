# cbkconvos
get your conversations in a tavern compatible way :)

## usage

1. install go
2. put your api token in the `bearer` file first line
3. run `go mod tidy`
4. run the command this way `go run . -c <yourConversationId> --charactername <yourCharacterName>`

you can also use `-h` to check what the params do and whats their format

example:
```
$ go run . -c c1ff34u3847382rd --charactername Megumin
```

another example with "your" user name set
```
$ go run . -c c1ff34u3847382rd --username Kazuma --charactername Megumin
```

This script creates two files:
- `backstory_167984154643401.txt` contains the gaslight + character info that you inputted in -- sometimes it doesn't work because the gaslight isn't sen't, but sometimes it does, cbk has no say on when it will give you a backstory and when it wont, so essentially this might just be the first message
- `167984154643401.jsonl` contains the actual messages in a format compatible with tavern
