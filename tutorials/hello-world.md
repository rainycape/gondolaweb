So now that you've installed Gondola, let's write our first app.

## Creating the app directory

Gondola comes with the *gondola* command, which can do a lot of useful things, but for now, we'll start by using its *new* subcommand to create our app directory. Start by typing this command from your GOPATH:

```bash
gondola new -template=blank myapp
```

This will create a directory named *myapp* containing a few files, essential for every Gondola application. Let's take a look at them:

- **app.conf**: This is the configuration file used for running the live app. Here we'll define settings like the database, the cache or other simpler stuff, like the port our app will listen on.
- **dev.conf**: This is the configuration file used during development. Note that the gondola development server uses this file by default, but will happily use *app.conf* if *dev.conf* does not exists. If your project uses the same configuration for development and production, you can safely remove this file.
- **assets**: This directory contains the static assets, like Javascript or CSS files that our app will use. Note that this directory is configurable, but it's usually better to keep the default name, since other people editing your code will inmediately recognize the name.
- **tmpl**: This directory contains the templates, which are used to generate HTML pages or outgoing emails. Again, this directory is also configurable, *tmpl* is just the default value.


## Our first handler

Now that we've seen the basic directory structure, let's write our first handler. Open your favorite text editor and type this code. Please, don't just copy and paste it, **type it**!

```go
package main

import (
        "gnd.la/app"
)

func MainHandler(ctx *app.Context) {
        ctx.WriteString("Hello world!")
}
```

Now save it as *handlers.go* in our application directory. Let's see what this code does.

- First, it imports **gnd.la/app**, which contains the central pieces of a Gondola application. It defines types like App and Context. The latter is specially important, since it will let us access the database, the cache, the blobstore and pretty much every subsystem of our application, and also send back responses to the client. Note that *app.Context* implements the [http.ResponseWriter](1), so whenever you see a function that accepts an *http.ResponseWriter*, you can pass an *app.Context*.

- Next, it defines the function named *MainHandler* which takes only one argument of type *\*app.Context* and has not return values. This is a very important type of function in Gondola, also known as an [app.Handler](2). Handlers are attached to endpoints and generate responses from requests (we'll see how to attach them in a bit).

- Finally, it calls WriteString() on the received Context. This just sends the string "Hello world!" to the client.

Now that we have our first handler, let's write an app and attach our handler to an endpoint.


## Tying it all together

Let's go back to our favorite text editor and type this code, saving it as *app.go*.

```go
package main

import (
        "gnd.la/app"
        "gnd.la/config"
)

var (
        App *app.App
)

func init() {
        config.MustParse()
        App = app.New()
        App.HandleNamed("^/$", MainHandler, "main")
}
```
And now let's examine what this code does.

- First, it imports *gnd.la/app*, as we've done before, but it also imports *gnd.la/config*. The *config* package handles all the configurable variables of a Gondola application, providing a few common settings and letting you add your own ones.

- It declares a global *App* variable without initializing it. This is important, since our App can't be initilized until the configuration has been parsed.

- Then we declare an *init()* function. Remember that *init()* functions are executed before main. This is where we're going
to parse the configuration file and initialize our App. It's important to parse the configuration first, since it will read our configuration
file and prepare some values for the App that will initialized afterwards. The purpose is two-fold: by preparing our App in *init()* we're sure the rest of the package won't be presented with an uninitializated App instance and, most importantly, our initialization code will also run in environments which don't allow a *main()* function, like Google App Engine.

- The first line in *init()* parses the configuration and stores the result in our global *Config* variable. Note that we use the function *config.MustParse()*, which has a Must prefix. This a commonly used idiom in Gondola for shorthand functions that panic rather than returning an error, so this means we don't have to explicitely check for errors here, the function will do it for us. For a complete reference of the *gnd.la/config* package, take a look a its [documentation](3).

- Once the configuration has been parsed, we can initialize our application. We'll always do this using the same function, *app.New()*, which takes no arguments and returns a new application, using the settings from the global configuration.

- Finally, we set up our handler to respond to requests at / using *App.HandleNamed()*. Note that the first parameter must be a valid regular expression, as understood by the [regexp](4) Go package. The second argument is our handler function, while the third argument sets the handler name, allowing us to obtain its URL from a template or handler using the *reversing* functions. Note that there's also a simpler version or *HandleNamed()*, called *Handle()* which doesn't take a name (and, thus, the handler URL can't be reversed).


There's just one little bit left, our *main()* function.
Go to your favorite text editor and type the following code into *main.go*:

```go
package main

func main() {
        App.MustListenAndServe()
}
```

This code is really simple, all it does is tell our app to start listening an accepting requests. Note that we do this in *main()* rather than on *init()* because we don't want the app to start listening when running tests or when running on App Engine (a different mechanism is used in that environment).


## Building and launching our first app

Gondola comes with a development server, which automatically compiles and rebuilds your code as necessary. To start it, type the following in your app directory (which, if you followed the example, will be called myapp):

```bash
gondola dev
```

This will print a few messages, like the port the development server is listening on (usually 8888), and start building your app. It should also automatically open your browser pointing to the development server (unless you're running this command in a remote shell, in that case you'll have to open it manually).

If your application can't be compiled due to some errors, the development server will show them to you in your browser. Just double-check that you typed the code correctly and save the files after changing them. The development server will detect the changes and try to compile them again.

Once your application builds and start, you'll see a very minimal web page which only says "Hello world!"

## Where to go next?

First, you can obtain the final code of this tutorial, with comments and a few more examples by creating a project with the hello template:

```bash
gondola new -template=hello hello
```

If you're interested in running your app in App Engine, it's recommended that you read the [Gondola on App Engine]({{ reverse_article "app-engine" }}) document before following the rest of the tutorials.

So you've already written your first Gondola application. Easy, wasn't it? Continue to the next tutorial: [Using Gondola templates]({{ reverse_article "gondola-templates" }}).

[1]: http://gondolaweb.com/doc/pkg/net/http#type-ResponseWriter
[2]: http://gondolaweb.com/doc/pkg/gnd.la/app#type-Handler
[3]: http://gondolaweb.com/doc/pkg/gnd.la/config
[4]: http://gondolaweb.com/doc/pkg/regexp


[title] = Hello World with Gondola
[slug] = hello-world
[synopsis] = This article includes instruction for installing Gondola from scratch.
[priority] = 2
