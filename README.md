### Trello API client in go

This is the Trello api part of a github/trello integration used at my workplace,
meant to enable team of developers to automatically include comments on trello cards.

This is by no means full featured Trello API, It contains the absolute minimum I need.
Having said that It should be fairly straightforward to implement other parts of trello API.

For documentation do

    godoc .

To run the tests you will need to create the test_config.json from the attached sample.
You will need following.

*  Trello API application key
*  Authorization token
*  A list to move the test card to
*  A test card id or the last part of the shortlink, card must have FOOBAR as its description.

Then just run it with

    go test -v
