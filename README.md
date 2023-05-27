<h2 align="center"><b>Box, Box! Server</b></h2>
<h4 align="center">A proxy written in Go for Box, Box!</h4>

[![GitHub releases](https://img.shields.io/github/release/BrightDV/BoxBox-Server?style=for-the-badge)](https://github.com/BrightDV/BoxBox-Server/releases/latest)
[![GitHub issues](https://img.shields.io/github/issues/BrightDV/BoxBox-Server?style=for-the-badge)](https://github.com/BrightDV/BoxBox-Server/issues)
[![GitHub forks](https://img.shields.io/github/forks/BrightDV/BoxBox-Server?style=for-the-badge)](https://github.com/BrightDV/BoxBox-Server/network)
[![GitHub stars](https://img.shields.io/github/stars/BrightDV/BoxBox-Server?style=for-the-badge)](https://github.com/BrightDV/BoxBox-Server/stargazers)
[![GitHub license](https://img.shields.io/github/license/BrightDV/BoxBox-Server?style=for-the-badge)](https://github.com/BrightDV/BoxBox-Server/blob/main/LICENSE)

## Installation
At first, you need Go installed.

- To install Go, follow the [installation instructions](https://go.dev/doc/install).

- To verify whereas Go is installed, run `go --version` in your terminal.

Step-by-step installation of Box, Box! server:
```
git clone https://github.com/BrightDV/BoxBox-Server.git
cd BoxBox-Server-main
go run main.go
```
That's it! Your Box, Box! server is now up and running on your device, on the port 8080 by default!

* If you wish to change the port used, open the `main.go` file with your favorite text editor and edit the line 43 by replacing `8080` with the port that you want to use.

* By default, Box, Box! Server is accessible from any website. If you want to only accept some websites to use your server, edit the line 42 of the `main.go` file and replace the `*` with your website(s) base URL(s) (e.g., the origin of the request).


## Box, Box! vs Box, Box! Server
At first, Box, Box! was an Android application. Then, thanks to Flutter, Box, Box! could be deployed as a website. Sadly, almost all the requests were blocked by the CORS of your browser... So here is Box, Box! Server: a proxy designed to provide you all the news and the results without being blocked.

> Note that Box, Box! Server is meant to be used with Box, Box! (see [the wiki]() to see the steps in order to use your server).

## License
[![GNU GPLv3 Image](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html)  

```
Box, Box! Server is Free Software: You can use, study, share, and improve it at
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU General Public License](https://www.gnu.org/licenses/gpl.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
```

## Contributions
Contributions are **very welcome**! I am only a newbie in Go so feel free to modify and improve the code, and, of course, open a PR ;)!