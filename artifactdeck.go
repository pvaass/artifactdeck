package artifactdeck

// Version Current ADC version
const Version = 2

// Prefix The prefix for the result string
const Prefix = "ADC"

const headerSize = 3

// Hero A hero in a deck
type Hero struct {
	ID int
	Turn int
}

// Card A card in a Deck
type Card struct {
	ID int
	Count int
}

// Deck The unencoded deck
type Deck struct {
	Name string
	Heroes []Hero
	Cards []Card
}


// MarshalText Encodes a Deck struct to a base64 encoded byte array
func (deck *Deck) MarshalText() ([]byte, error) {
	text, err := encode(*deck)

	return []byte(text), err
}


// UnmarshalText Decodes a base64 encoded byte array to a Deck struct
func (deck *Deck) UnmarshalText(text []byte) error {
	decoded, err := decode(string(text))

	if err != nil {
		return err
	}

	*deck = decoded	

	return nil
}