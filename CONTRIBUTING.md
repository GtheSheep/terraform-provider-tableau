# Contributing

## Running Unit Tests
See `env.vars.example` for an example of the set of environment variables that need to be set in order to run acceptance tests.  
In order to run Tableau Server (not Cloud) specific tests, an additional `TF_ACC_SERVER` environment variable must be set, as some API methods do not apply to Tableau Cloud, i.e. Sites.  
