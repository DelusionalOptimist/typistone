# Typistone
[typistone](https://github.com/DelusionalOptimist/typistone) is a TUI game that lets you practice typing drills and compete with other players across the world.

### Try it out

#### Setup
* First deploy [typistone-server](https://github.com/DelusionalOptimist/typistone-server) (for now).
* Then build and run the CLI
```
go build -o ./typistone
```

#### Play
* Play singleplayer:
```
./typistone singleplayer
```
* Play multiplayer:
```
./typistone multiplayer create --lobby-size <lobby_size>
```
* Join a multiplayer game:
```
./typistone multiplayer join --lobby-id <lobby_id>
```
