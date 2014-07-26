Gondola is written in Go, so in order to install Gondola
a working Go installation is required.

## Installing Go

To install Go, follow the [installation instructions](http://golang.org/doc/install) 
on the Go website.

If you're new to Go, it's also recommended to follow the
[Go Tour](http://tour.golang.org) to familiarize yourself
with the language.

## Installing required tools

Gondola is published as a *git* repository, while some of
its dependencies are made available with *mercurial*. In
order to download Gondola, you'll need both git and mercurial
installed on your computer. Please, refer to your operating
system section to learn how to install these tools.


### Windows

Download the git installer from the [git Windows downloads](http://git-scm.com/download/win)
and run the installer. Then go to [mercurial downloads](http://mercurial.selenic.com/downloads)
and choose the Mercurial Inno Setup installer installer for either 32 or 64 bits. If you're unsure about
which version you need, download the 32 bits one. Run its installer and you're set.

### Mac OS X

In OS X, git can be installed either by installing [Xcode](https://developer.apple.com/xcode/downloads/) and then
downloading the Command Line Tools from the Xcode -> Preferences -> Downloads menu or by downloading an installer from
the [git OS X downloads](http://git-scm.com/download/mac). To install mercurial, download the latest OS X package
from [mercurial downloads](http://mercurial.selenic.com/downloads).

### Linux

On Debian based distributions, like Ubuntu, open a terminal and type the following:

```sh
sudo apt-get install git mercurial
```

On Red Hat based ones, like Fedora, Centos or Suse, type this instead.

```sh
sudo yum install git mercurial
```

If you're running any other type of distro, you probably know how to install these tools yourself.

## Installing Gondola

Once Go is installed, open a command prompt (if you're on
Windows press Win+R, type cmd and press Enter) and type
the following:

```sh
go get -v gnd.la/...
```

Please note that this will take a few minutes to complete and will download several packages. Some
non-essential packages might fail to install due to required depedencies, but that's usually not
a problem. Now type the following command:

```sh
go install gnd.la/cmd/gondola
```

If it completes without errors, Gondola is most likely correctly installed. Check the last step
in this tutorial to make sure.

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
with the [Hello World]({{ reverse_article "hello-world" }}) tutorial.


[title] = Installing Gondola
[synopsis] = Read this tutorial after installing Gondola.
[updated] = 2014-07-26 06:27:40
[updated] = 2014-07-26 16:43:30
[priority] = 1
