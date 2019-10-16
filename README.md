## Usage
`go get github.com/gregmus2/ants-pkg` in your algorithm dir and execute
`go build -buildmode=plugin -o {OUT_FILE} {GO_FILE}` to build your plugin

Variable, which you want to use as algorithm have to implement Algorithm interface
and have your nickname