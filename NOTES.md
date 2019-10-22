# Development Notes

## TODO

This project is in active development. There is much work to be done.

- documentation
- error handling
- verbose and quiet mode
- testing
- terraform generators
- notification targets and expectations
- update templates to be generic

## Datadog Specific Features

This project is a cleaned up version of a shell script and simple go app. 
It was designed for internal use at datadog and thus makes several assumptions
about the state of things.

- slack expecation (should just be notificaiton targets)
- pager expectation (part of notifications)
- datacenters (part of a general substitution mechanism)


