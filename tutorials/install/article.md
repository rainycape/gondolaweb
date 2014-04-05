Gondola is written in Go, so in order to install Gondola
a working Go installation is required.

## Installing Go

To install Go, follow the [installation instructions](http://golang.org/doc/install) 
on the Go website.

If you're new to Go, it's also recommended to follow the
[Go Tour](http://tour.golang.org) to familiarize yourself
with the language.

## Installing Gondola

Once Go is installed, open a command prompt (if you're on
Windows press Win+R, type cmd and press Enter) and type
the following:

```sh
go get gnd.la/...
```

## Verify Gondola was installed

To make sure Gondola was installed correctly, type the
following in a command prompt:

```sh
gondola help
```

It should print the the help the Gondola command, which
we'll use in the next tutorial. If the command fails or
the gondola command can't be found, please, read the
previous steps carefully and make sure you did everything
right.

## Done!

Now that your Gondola installation is ready, continue
with the [Hello World]({{ reverse_article @Ctx "hello-world" }}) tutorial.
