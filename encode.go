package artifactdeck

import "sort"
import "errors"
import "html"
import "encoding/base64"
import "strings"

func encode(deck Deck) (string, error) {

	bytes, err := encodeDeckToBytes(deck)
	if err != nil {
		return "", err
	}

	deckCode, err := encodeBytesToString(bytes)
	if err != nil {
		return "", err
	}

	return deckCode, nil
}

func encodeDeckToBytes(deck Deck) ([]byte, error) {

	sort.Slice(deck.Heroes, func(i, j int) bool {
		return deck.Heroes[i].ID < deck.Heroes[j].ID
	})
	sort.Slice(deck.Cards, func(i, j int) bool {
		return deck.Cards[i].ID < deck.Cards[j].ID
	})

	heroCount := len(deck.Heroes)

	bytes := []byte{}
	version := Version<<4 | extractNBitsWithCarry(heroCount, 3)

	ok, bytes := addByte(bytes, version)

	if !ok {
		return nil, errors.New("Incorrect version")
	}

	dummyChecksum := 0
	checksumByte := len(bytes)

	ok, bytes = addByte(bytes, dummyChecksum)
	if !ok {
		return nil, errors.New("Can't append dummy checksum")
	}

	nameLength := 0
	if len(deck.Name) > 0 {
		sanitizedName := html.EscapeString(deck.Name)
		byteString := []byte(sanitizedName)

		if len(byteString) > 63 {
			deck.Name = string(byteString[:63])
		} else {
			deck.Name = string(byteString)
		}

		nameLength = len(deck.Name)
	}

	ok, bytes = addByte(bytes, nameLength)
	if !ok {
		return nil, errors.New("Can't append length of deck name")
	}

	bytes, err := addRemainingNumberToBuffer(bytes, heroCount, 3)
	if err != nil {
		return nil, errors.New("Can't append length of deck name")
	}

	prevCardID := 0
	for currentHero := 0; currentHero < heroCount; currentHero++ {
		hero := deck.Heroes[currentHero]
		if hero.Turn == 0 {
			return nil, errors.New("Can't have hero with turn 0")
		}

		bytes, err = addCardToBuffer(bytes, hero.Turn, hero.ID-prevCardID)
		if err != nil {
			return nil, errors.New("Can't add card to buffer")
		}

		prevCardID = hero.ID
	}

	prevCardID = 0

	for currentCard := 0; currentCard < len(deck.Cards); currentCard++ {
		card := deck.Cards[currentCard]
		if card.Count == 0 {
			return nil, errors.New("Can't have card with count 0")
		}
		if card.ID <= 0 {
			return nil, errors.New("Can't have card with id < 0")
		}

		bytes, err = addCardToBuffer(bytes, card.Count, card.ID-prevCardID)

		if err != nil {
			return nil, errors.New("Error while adding card to buffer")
		}
		
		prevCardID = card.ID
	}

	preStringByteCount := len(bytes)

	nameBytes := []byte(deck.Name)
	for currentNameByte := 0; currentNameByte < len(nameBytes); currentNameByte++ {
		ok, bytes = addByte(bytes, int(nameBytes[currentNameByte]))
		if !ok {
			return nil, errors.New("Error while appending name to bytestring")
		}
	}

	fullChecksum := computeChecksum(bytes, preStringByteCount-headerSize)
	smallChecksum := fullChecksum & 0x0FF

	bytes[checksumByte] = byte(smallChecksum)

	return bytes, nil
}

func encodeBytesToString(bytes []byte) (string, error) {

	if len(bytes) == 0 {
		return "", errors.New("Can't encode bytes of length 0")
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	encoded = strings.Replace(encoded, "/", "-", -1)
	encoded = strings.Replace(encoded, "=", "_", -1)

	return Prefix + encoded, nil
}
