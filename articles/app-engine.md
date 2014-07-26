## General considerations

Gondola is written to abstract you, the developer, from the underlying platform
where your app is running. This means that the same code will run on App Engine
and on bare metal servers with no changes in most cases.

However, due to the way App Engine works, you must avoid the following packages
in the Go standard library:

- net/http, use [gnd.la/net/httpclient](/doc/pkg/gnd.la/net/httpclient) instead.
- database/sql, use [gnd.la/orm](/doc/pkg/gnd.la/orm) instead.

Also, when possible, avoid using any App Engine APIs directly. Most subsystems in
Gondola are pluggable and have specific backends for App Engine. For example,
gnd.la/cache and gnd.la/blobstore can transparently work using the App Engine
backends when running on GAE and use any other backend anywere else (e.g. redis
for the cache or S3 for storing files).

## Installation

To develop and run Gondola applications in App Engine, you need to install Gondola
as detailed in the [Gondola Installation](/tutorials/installing-gondola/). Then,
install the Go [App Engine SDK](https://developers.google.com/appengine/downloads)
for your operating system.

## Development

Due to the restrictions imposed by the App Engine development environment, the
usual *gondola dev* command won't work correctly. Instead, users must use the
**gondola gae-dev** command.

## Testing

To run tests in your application, use the **gondola gae-test** command. For more
information on Gondola tests, see [gnd.la/app/tester](/doc/pkg/gnd.la/app/tester).


## Deployment

In order to deploy a Gondola application to App Engine, all its assets must be
pregenerated (usually, these are generated on-demand). To do so, just use the
**gondola gae-deploy** command. It will generate all the required assets and then
upload your app to App Engine using the GAE SDK itself.


[title] = Gondola On App Engine
[synopsis] = Learn how to develop and deploy Gondola applications on Google App Engine.
