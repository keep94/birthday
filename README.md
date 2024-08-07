# birthday

A birthday and special day reminder system.

This program is a webserver that reads names and birthdays from a text file and shows upcoming special days for those people. A special day is any of the following:

- A person's birthday
- A person turning a multiple of 1000 days old
- A person turning a multiple of 100 weeks old
- A person turning a multiple of 100 months old

## Format of the text file

The text file is a TSV file (tab delimited file) with two fields per record:

- Name
- Date of Birth (MM/dd/yyyy format. Use MM/dd format if you don't know the year a person was born)

A sample text file may look like this:

```
John Smith	3/25/1967
Bill Shaw	12/07/1973
Katie Long	3/21/2010
Merna Heitcamp	5/17
```

Although the records in the file can be in any order, I recommend ordering by name to make it easier to make updates to the file.

## Building

To build the server, do the following:

- Download golang
- In your home directory, create a go/src/github.com/keep94 folder
- cd to that folder.
- Then run `git clone git@github.com:keep94/birthday.git`
- Run `cd birthday`
- Run `go install ./...` there
- You will find the executable at $HOME/go/bin/remind

## Running the server

Use `$HOME/go/bin/remind -file path/to/tsv/file -http ":8283"`

This tells the web server to use path/to/tsv/file for the birthdays and to
listen on port 8283. 

If you point your browser to `http://localhost:8283`, you will see the upcoming special events. You get redirected to `http://localhost:8283/home`

You can use `$HOME/go/bin/remind -file path/to/tsv/file` and the port defaults to 8080.

You can only see the first 100 special events.

The rest of this document assumes the webserver is listening on port 8080.

## Special tricks for viewing upcoming special days

By default, you see upcoming reminders for special days from today up to but not including 21 days from now. The reminders for today come first and are in italics.

### Want to see if any special days happened yesterday or the day before

Point your browser to `http://localhost:8080/home?date=5/1` Where date is the month and day of the current year. If you want to go back to a prior year, you can use `http://localhost:8080/home?date=12/28/2023`

### Want to see special days for one person

Point your browser to `http://localhost:8080/home?q=perez&days=365` This shows only people with perez in their name and shows all special days up to but not including 365 days from now.

### Want to see only birthdays and no other special days

Point your browser to `http://localhost:8080/home?p=y`

The p parameter controls what types of special days show up. Special day types are as follows:

| Letter | Description |
| ------ | ----------- |
| y | traditional birthday |
| d | 1000 day multiple |
| w | 100 week multiple |
| m | 100 month multiple |
| h | 6 month multiple. Traditional birthdays and half birthdays |

If you wanted to see only traditional birthdays and 100 month multiples, you would use p=ym.

