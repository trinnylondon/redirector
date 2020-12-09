# About
This is a traefik pilot plug-in that will be used to do client redirection. 

# Design
The redirection is performed based on the following rules for requests to `/`:

1. If the incoming request has a cookie with a name in the `redirectCookies` list and a value found in `redirectionMap`, then the client gets a `301` redirect to that path.
2. If the incoming request has a cookie with a name in the `redirectHeaders` list and a value found in `redirectionMap`, then the client gets a `301` redirect to that path.
3. If none of the above rules are applied, the client gets a redirect to `defaultPath`.

If the request is for any path other than `/` then this plugin has no effect.