/*
Package starling provides a client for using the Starling API.

Usage:

	import "github.com/lildude/starling"

Construct a new Starling client, then call various methods on the API to access
different functions of the Starling API. For example:

	client := starling.NewClient(nil)

	// retrieve transactions for the current user
	txns, _, err := client.Transactions(ctx, nil)

The majority of the API calls will require you to pass in an access token:

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TOKEN"},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := starling.NewClient(tc)

The Starling API documentation is available at https://developer.starlingbank.com/docs.

*/
package starling
