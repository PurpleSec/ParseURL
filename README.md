# parseurl

Parseurl is a simple wrapper that fixes some of the weird issues that the standard Golang 'url.Parse' function does.

Fixes things such as

- "localhost:8080"
  url.Parse: Host == "", Scheme == "localhost:8080"
- "10.10.10.10/url/"
  url.Parse: Returns and error <why?>
- "localhost:"
  url.Parse: Host == "localhost:" <doesn't strip the ':'>

This package adds in checks for invalid values returned by 'url.Parse' such as the Host field being empty. All non-standard errors wrap the error "ErrInvalidURL" to assist in indication of the error.

This library is a drop-in replacement for the "url.Parse" function. Just import "github.com/PurpleSec/parseurl" and go!
