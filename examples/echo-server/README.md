## echo server

```
$ go build
$ ./echo-server
2021/11/14 16:16:27 server addr => 127.0.0.1:54456
```


## client

Use [netcat](https://en.wikipedia.org/wiki/Netcat)

Please enter `^D` (control + D) after you enter any characters.

```
$ nc 127.0.0.1 54456
hello
 _______
< hello >
 -------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```
