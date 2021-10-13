# valor

## Description

Valor is a tool that can be used to validate and evaluate resource.

## Quick Start

### 1. Clone Repo

Clone this repository into your local by running the following command:

```zsh
git clone github.com/gojek/optimus-extension-valor
```

### 2. Create Recipe

Create a recipe named `valor.yaml` or rename the recipe file in
this project directory from `valor.example.yaml` into `valor.yaml`.
The content of the recipe should look like the following.
Assets referred by the specified recipe are available under `./example` directory.

```yaml
resources: 
- name: user_account
  type: file
  path: ./example/resource/user_account.json
  format: yaml
  framework_names:
  - user_account_validation

frameworks:
- name: user_account_validation
  schemas:
  - name: user_account_rule
    type: file
    format: json
    path: ./example/schema/user_account_rule.json
    output_is_error: true
  procedures:
  - name: enrich_user_account
    type: file
    format: jsonnet
    path: ./example/procedure/enrich_user_account.jsonnet
  output_targets:
  - name: std_output
    type: std
    format: json
```

In this example, Resource named `user_account` will be evaluated agaist
framework named `user_account_validation`. The targetted framework
specifies three things:

* Schema follows JSON schema format to validate the Resource. It is
required to validate whether the structure for the `user_account`
Resource contains error or not.
* Procedure follows JSONNET format to evaluate the `user_account` Resource.
It can be for validation or any execution that can't be handled by
using Schema.
* Output Target specifies on how to write the output of either Schema or
Procedure or both.

### 3. Execute Pipeline

After the preparation is done, try running valor by executing the following
command:

```zsh
./out/valor
```

Or, if you have not build it, try building it by following [#HowToBuild](#how-to-build).
Make sure to have the require dependencies specified under [#Dependency](#dependency).
The output of the above solution should look like the following:

```zsh
user_account: ./example/resource/user_account.json
{
   "email": "valor@github.com",
   "is_active": true,
   "is_valid": true,
   "membership": "premium"
}
```

## Dependency

### GO Language

The GO programming language version `go1.17.1` need to be installed
in the system. Go to [this link](https://golang.org/doc/install) and follow
the instruction to install. To check GO version on the environment, run the following command:

```bash
go version
```

example output:

```bash
go version go1.17.1 darwin/amd64
```

## How to Test

In this project root directory, run the following command to test:

```bash
make test
```

To show a more complete coverage and uncovered lines:

```bash
make coverage
```

You can check into `coverage.html` file in root project directory.
This command also will open interactive coverage tool in your browser if you have one.

## How to Build

To build the binary executable, in this project root directory, run the following command:

```bash
make bin
```

There will be a new directory named `out` with an executable file `valor` as the result of the built project.

## How to Run

In order to run this project, after building the binary executable, run
the following command in this project root directory:

```bash
./out/valor
```
