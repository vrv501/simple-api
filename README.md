# go-template
Starter Repo Config for Building Go Apps in a debian amd64 machine

## Tools
- Go Compiler: [Download-Link](https://go.dev/doc/install)

- Install pre-requisite tools  
  ```bash
  $ sudo apt install -y gcc make curl jq
  ```

- Add GOBIN to system PATH
  ```bash
  $ export PATH=$PATH:$(go env GOPATH)/bin
  ```

- Install golangci-lint  
  ```bash
  $ tag_name=$(curl -s -L -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" https://api.github.com/repos/golangci/golangci-lint/releases/latest | jq -r .tag_name)
  $ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin "$tag_name"
  ```

- Install mockgen
  ```bash
  $ tag_name=$(curl -s -L -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" https://api.github.com/repos/uber-go/mock/releases/latest | jq -r .tag_name)
  $ go install "go.uber.org/mock/mockgen@$tag_name"
  ```
- Add validation of pet images size
- Check for get filters atleast one parameter is provided
