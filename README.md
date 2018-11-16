ArtifactDeckCode port for go.

# Installing
`go get github.com/pvaass/artifactdeck`

# Usage
```go
import artifact "github.com/pvaass/artifactdeck"

func main() {
  var deckCode string
  var deck artifact.Deck

  // Create Deck from string
  err = deck.UnmarshalText(deckCode)
  if err != nil {
    panic(err)
  }

  // Create string from Deck
  text, err = deck.MarshalText()
  if err != nil {
    panic(err)
  }
}

# Docs
https://godoc.org/github.com/pvaass/artifactdeck