# Design Log

## 2024/11/15

Initially the application will be a single-user CLI program. User data will be stored in a local
SQLite database. The initial usecase is to allow the user of the application to enter the amount
of a medication that they have, and the rate at which the medication is used, and the program will
calculate when the next refill is required. While associating medications with the person or animal
that takes them would be nice, it is not necessary for the first implementation. While it might be
good to plan ahead for that, I will intentionally not do that here, to introduce some more difficult
upgrade and compatability scenarios later on.

The CLI program will be called `rxm` for "Rx Manager" (Prescription manager.) This should be 
quick and easy to type and remember. The CLI program will have several subcommands to manipulate the
information in the database.

* `rxm list` -- This will list all of the medications entered into the system so far, along with the
  expected current amount of the medication remaining, and the expected next refill date.
* `rxm add [name] [amount] [rate]` -- This will add a new medication to the database.
* `rxm rm [name]` -- This will remove a medication from the database.
* `rxm up [name] quantity [amount]` -- This will update the quantity of a medication in the
  database.
* `rxm up [name] rate [rate]` -- This will update the burn rate of a medication in the database.
* `rxm show [name]` -- This will show detailed information about a single medication.

A `-config` flag can be added to `rxm` to allow the location of configuration directory to be set.
If this flag is not provided, the default location will be used. On Linux/macOS the default should
be in `~/.config/rxm/`, and on Windows it should be where the `APPDATA` environment variable points.

This is a very minimal but useful initial version of the application to implement. Now that we have
decided how the application will work from the user's perspective, we can now decide how to model
the information. This data model may indeed be simple enough to just store as a JSON or YAML file in
the configuration directory, but we've already decided to use sqlite3, so need to think of it in
terms of tables.

The above application is so simple, the data model basically writes itself. A `prescription` has a
unique `name`, a quantity, and a burn rate. The quantities can be represented using `DECIMAL(10, 2)`
allowing for precision up to 1/100th of one unit of medication, which should be enough precision.
The medication names can be fairly long, but not absurdly, so a 64-character limit should be
sufficient, therefore a `VARCHAR(64)` will work. These pieces of information are sufficient to
determine the next refill date _now_, but they are not sufficient to keep that date stable over
time. Therefore we also need a `DATE` for when the quantity was last updated. Putting this all
together, we really only need a single table with four fields.

```sql
CREATE TABLE prescription (
    prescription_id VARCHAR(64) PRIMARY KEY,
    quantity DECIMAL(10, 2),
    rate DECIMAL(10, 2),
    updated_at DATE
);
```

Note that this design loses quantity-over-time information. This is not a great design choice, but
is being made intentionally here so I'll have to dig myself out of it later. Might be better to have
a table of quantity updates, so it can be tracked over time. For the initial implementation a single
table should be fine.

## 2024/11/16

Took the existing CLI, and created an HTTP JSON REST API of essentially the same shape. A web
frontend can be built over this API. It might be possible to allow the UI/API to be launched by the
CLI, for example via `rxm ui` and it would run the backend server in the foreground and also launch
the user's default browser to the URL needed to hit the UI.