# cli-mahjong

## setup

golang v1.13.12

download golang code from https://golang.org/dl/

```bash
$ go version
go version go1.13.12 darwin/amd64
```

## run

activate the game sever. default running port is 8080

```bash
go run main.go
```

connecting from clients, at least four players required to start a match.

```bash
netcat localhost 8080
```

## playing

this is game screen. the player having `>>` is the turn player and can discard a tile by typing tile's name

```
                            ┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐
                            │  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  │
                            └──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘
┌──┐
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘                                                                        ┌──┐
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
└──┘                                                                        │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            └──┘

┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐    ┌──┐
│m1││m3││m4││m8││m8││m9││p3││p5││p6││p8││p9││s1││西│    │s3│
└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘    └──┘

>>
```

you can take huro actions (chii, pon, kan) in other player's turn. 
chooing huro action can be done by typing action name (chii, pon, kan), and index. 
in the below example, you would make s7,s8, and s6 meld as your tiles


```bash
                            ┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐
                            │  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  ││  │
                            └──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘
┌──┐
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘
│  │
└──┘                                            ┌──┐
│  │                                            │中│
└──┘                    ┌──┐                    └──┘                        ┌──┐
│  │                    │s6│                    └──┘                        │  │
└──┘                    └──┘                                                └──┘
│  │                    └──┘                                                │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                                        └──┘
│  │                                                                        │  │
└──┘                                                ┌──┐                    └──┘
│  │                                                │南│                    │  │
└──┘                        ┌──┐                    └──┘                    └──┘
└──┘                        │発│                    └──┘                    │  │
                            └──┘                                            └──┘
                            └──┘                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            │  │
                                                                            └──┘
                                                                            └──┘
                                                                    ┌──┐┌──┐┌──┐
┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐                            │西││西││西│
│m5││m6││p2││s3││s4││s5││s6││s7││s8││s8│                            └──┘└──┘└──┘
└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘                            └──┘└──┘└──┘

do you do Chii 0: s4,s5   1: s5,s7   2: s7,s8
chii2
```


## naming

https://en.wikipedia.org/wiki/Japanese_Mahjong
