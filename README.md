# go-medminder

This repository contains a simple utility for managing and tracking prescriptions. It is a work in
progress, started primarily to learn the Go programming language.

Presently the only thing implemented is a single-user CLI program, but this can be expanded to a
multi-user web application in the future.

## Usage

Build and install like any Go package. The CLI program is called `rxm` (short for Prescription
Manager). There are several subcommands:

- `rxm add [name] [quantity] [rate]` -- Adds a new prescription named `name` with `quantity` amount
  of the medication, where `rate` amount of medication is used per day.
- `rxm rm [name]` -- Remove the prescription named `name`. Permanently loses information.
- `rxm up [name] quantity [quantity]` -- Update the current `quantity` for the prescription named
  `name`.
- `rxm up [name] rate [rate]` -- Update the current `rate` for the prescription named `name`.
- `rxm ls` -- List basic information about all prescriptions in the database.
- `rxm ls [name]` -- List detailed information about the prescription named `name`.

On Linux, by default, the database file will be placed at `$XDG_CONFIG_DIR/go-medminder/db.sqlite3`,
or if `XDG_CONFIG_DIR` is not set, then `~/.config/go-medminder/db.sqlite3` is the default. On
Windows `$APPDATA/go-medminder/db.sqlite3` is the location. A flag can be added to the `rxm`
invocation to override this default location: `-db [path]` where `path` is the absolute or relative
path to the desired location of the database file. If the database file does not already exist, it
is created upon invocation of `rxm`.