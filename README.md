# Rent-N-Go Backend
A complete backend service that empowers Rent-N-Go. This service not only provide API for 
the end user, but also provide admin panel to help maintain the end user.

Written in Go.

# Compiling Rent-N-Go

## Prerequisites 

- GNU Make
- NodeJS v16 or above
- PNPM v7 or above
- Go v1.20 and above

## Instructions

For Unix Based operating system, this project has provided `Makefile` which will take care the 
installation automatically. The `Makefile` will assume that you have `pnpm` and `go` binaries
installed.

Simply run `make` and you're set.

If you don't want to use `pnpm` you can opt it out to use npm instead by running `make install_npm` and 
you're set.

# Windows Users

Rent-N-Go does not provide official way to build compilation environment for Windows user though. But windows
user can still compile this adding some prerequisites:

- [Git's bash](https://git-scm.com/downloads) (we only need bash in here. You can switch to alternatives like [Cygwin](https://www.cygwin.com/) if you want to)
- [Choco Package Manager](https://chocolatey.org/) (we need this in order to install `make` for windows)
- [GNU Make](https://community.chocolatey.org/packages/make) (in order to run make command)

After installing choco make sure to add it in path so Git Bash can access it too. From there, type `choco install make`
to perform installation for Gnu's Make.

> Alternatively, you can see the manual install below.

## Manual Install

`$PKG` can either be `npm` or `pnpm`.

- `$PKG` install
- `go mod download`
- `go run "./generator/*.go"`
- `go run *.go migrate`
- `go run *.go seed`

And you're set.

# Configuration

## Prerequisites

- MySQL v8 or above

## Steps

- Copy `.env.example` to `.env`. 
- Edit configuration as you need
- Create a new database in your MySQL that match your `DB_NAME`
- Simply run the project by performing `make run` and you're set.

# Code Generation

Rent-N-Go takes an advantages and fully relies on Gorm Gen. A code generation package provided
by Gorm. Whenever you'll need to regenerate a code, simply run `make generate`.