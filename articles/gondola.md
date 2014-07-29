Gondola provides a command line tool named **gondola**. This article
lists all its options and subcommands.
## Global flags

 - **q** *\(bool\)*: Disable verbose output *default: false*


## Commands


- ## dev

    Start the Gondola development server



    Flags:

     - **dir** *\(string\)*: Project directory *default: \.*
     - **port** *\(int\)*: Port to listen on *default: 8888*
     - **config** *\(string\)*: Configuration file\. If empty, dev\.conf and app\.conf are tried in that order
     - **tags** *\(string\)*: Build tags to pass to the Go compiler
     - **no\-debug** *\(bool\)*: Disable AppDebug, TemplateDebug and LogDebug \- see gnd\.la/config for details *default: false*
     - **no\-cache** *\(bool\)*: Disables the cache when running the project *default: false*
     - **profile** *\(bool\)*: Compiles and runs the project with profiling enabled *default: false*
     - **race** *\(bool\)*: Enable \-race when building\. If the platform does not support \-race, this option is ignored *default: false*
     - **no\-browser** *\(bool\)*: Don't open the default browser when starting the development server *default: false*
     - **v** *\(bool\)*: Enable verbose output *default: false*



- ## new

    Create a new Gondola project

    Usage: ```gondola new <dir>```

    Flags:

     - **template** *\(string\)*: Project template to use *default: hello*
     - **list** *\(bool\)*: List available project templates *default: false*
     - **gae** *\(bool\)*: Create an App Engine hybrid project *default: false*



- ## build

    Build packages

    Usage: ```gondola build [package-1] [package-2] ... [package-n]```

    gondola build is a wrapper around go build\.
    
    Before building a package, gondola build checks that all its dependencies
    exist and automatically downloads the missing ones\.
    
    The \-go option can be used to set the go command that will be called\. All the
    remaining options are passed to go build unchanged\.

    Flags:

     - **go** *\(string\)*: Command to run the go tool *default: go*
     - **race** *\(bool\)*: Enable data race detection *default: false*
     - **x** *\(bool\)*: Print the commands *default: false*
     - **v** *\(bool\)*: Print the names of packages as they are compiled *default: false*
     - **ccflags** *\(string\)*: Arguments to pass on each 5c, 6c, or 8c compiler invocation
     - **compiler** *\(string\)*: Name of compiler to use, as in runtime\.Compiler \(gccgo or gc\)
     - **gccgoflags** *\(string\)*: Arguments to pass on each gccgo compiler/linker invocation
     - **gcflags** *\(string\)*: Arguments to pass on each 5g, 6g, or 8g compiler invocation
     - **ldflags** *\(string\)*: Arguments to pass on each 5l, 6l, or 8l linker invocation
     - **tags** *\(string\)*: A list of build tags to consider satisfied during the build



- ## clean

    Cleans any Gondola packages which use conditional compilation \- DO THIS BEFORE BUILDING A BINARY FOR DEPLOYMENT \- see golang\.org/issue/3172





- ## profile

    Show profiling information for a remote server running a Gondola app

    Usage: ```gondola profile <url>```

    Flags:

     - **method** *\(string\)*: HTTP method *default: GET*
     - **data** *\(string\)*: Optional data to be sent in the request in the form k1=v1&k2=v2\.\.\.



- ## gen\-app

    Generate boilerplate code for a Gondola app from the appfile\.yaml file



    Flags:

     - **release** *\(bool\)*: Generate release files, otherwise development files are generated *default: false*



- ## bake

    Converts all assets in \<dir\> into Go code and generates a VFS named with \<name\>

    Usage: ```gondola bake -dir=<dir> -name=<name> ... additional flags```

    Flags:

     - **dir** *\(string\)*: Root directory with the files to bake
     - **name** *\(string\)*: Variable name of the generated VFS
     - **o** *\(string\)*: Output filename\. If empty, output is printed to standard output
     - **ext** *\(string\)*: Additional extensions \(besides html, css and js\) to include, separated by commas



- ## random\-string

    Generates a random string suitable for use as the app secret



    Flags:

     - **length** *\(int\)*: Length of the generated random string *default: 64*



- ## rm\-gen

    Remove Gondola generated files \(identified by \*\.gen\.\*\)

    Usage: ```gondola rm-gen [dir]```



- ## make\-messages

    Generate strings files from the current package \(including its non\-package subdirectories, like templates\)



    Flags:

     - **o** *\(string\)*: Output filename\. If empty, messages are printed to stdout\. *default: \_messages/messages\.pot*



- ## compile\-messages

    Compiles all po files from the current directory and its subdirectories



    Flags:

     - **o** *\(string\)*: Output filename\. Can't be empty\. *default: messages\.go*
     - **ctx** *\(string\)*: Default context for messages without it\.



- ## gen

    Perform code generation in the current directory according the rules in the config file



    Flags:

     - **genfile** *\(string\)*: Code generation configuration file *default: genfile\.yaml*



- ## gae\-dev

    Start the Gondola App Engine development server





- ## gae\-test

    Start serving your app on localhost and run gnd\.la/app/tester tests against it



    Flags:

     - **v** *\(bool\)*: Enable verbose tests *default: false*



- ## gae\-deploy

    Deploy your application to App Engine








[title] = The Gondola command line tool
[synopsis] = Help for the gondola command.
[updated] = 2014-07-27
