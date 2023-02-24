# Rent-N-Go Backend
A complete backend service that empowers Rent-N-Go. This service not only provide API for 
the end user, but also provide admin panel to help maintain the end user.

Written in Go.

# Installation

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

For windows users that not using `bash`. You must run the `make` command with `windows_` prefix. Here's
a complete lists and their equivalent:

- `make windows` -> `make`
- `make windows_install_npm` -> `make install_npm`
- `make windows_run` -> `make run`

> Other command should work in both OS.

## Manual Install

`$PKG` can either be `npm` or `pnpm`.

- `$PKG` install
- `go mod download`
- `go run "./generator/*.go"`

And you're set.

# Configuration

## Prerequisites

- MySQL v8 or above

## Steps

- Copy `.env.example` to `.env`. 
- Edit configuration as you need
- Create a new database in your MySQL that match your `DB_NAME`
- Simply run the project by performing `make run` and you're set.