package artifactdeck

import "errors"
import "encoding/base64"
import "strings"
import "html"

func decode(deckCode string) (Deck, error) {
	deckBytes, err := deckStringToBytes(deckCode)

	if err != nil {
		return Deck{}, err
	}

	deck, err := deckBytesToStruct(deckBytes)


	return deck, err
}


func deckStringToBytes(deckCode string) ([]byte, error){
	
	if !strings.HasPrefix(deckCode, Prefix) {
		return nil, errors.New("No prefix on deck code")
	}

	deckCode = strings.TrimPrefix(deckCode, Prefix)

	deckCode = strings.Replace(deckCode, "-", "/", -1)
	deckCode = strings.Replace(deckCode, "_", "=", -1)


	bytes, err := base64.StdEncoding.DecodeString(deckCode)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}


func deckBytesToStruct(bytes []byte) (Deck, error) {

	currentByteIndex := 0
	totalBytes := len(bytes)

	versionAndHeroes := int(bytes[currentByteIndex])
	currentByteIndex++
	version := versionAndHeroes >> 4

	if Version != version && version != 1 {
		return Deck{}, errors.New("Version mismatch")
	}

	checksum := int(bytes[currentByteIndex])
	currentByteIndex++


	stringLength := 0
	if version > 1 {
		stringLength = int(bytes[currentByteIndex])
		currentByteIndex++
	}
	totalCardBytes := totalBytes - stringLength



	computedChecksum := 0
	
	for i := currentByteIndex; i < totalCardBytes; i++ {
		computedChecksum += int(bytes[i])
	}
	masked := computedChecksum & 0xFF

	if checksum != masked {
		return Deck{}, errors.New("Checksum error")
	}
	
	heroCount, err := readVarEncodedUint32(versionAndHeroes, 3, bytes, &currentByteIndex, totalCardBytes)

	if err != nil {
		return Deck{}, nil
	}

	var heroes []Hero

	prevCardBase := 0
	for currHero := 0; currHero < int(heroCount); currHero++ {
		card, err := readSerializedCard(bytes, &currentByteIndex, totalCardBytes, &prevCardBase)

		if err != nil {
			return Deck{}, nil
		}

		heroes = append(heroes, Hero{
			ID: card.ID,
			Turn: card.Count,
		})
	}


	var cards []Card
	prevCardBase = 0
	for currentByteIndex < totalCardBytes {
		card, err := readSerializedCard(bytes, &currentByteIndex, totalCardBytes, &prevCardBase)

		if err != nil {
			return Deck{}, nil
		}
		cards = append(cards, card)
	}

	name := ""
	if currentByteIndex <= totalBytes {
		nameBytes := bytes[(totalBytes - stringLength):]
		name = string(nameBytes)
	}

	return Deck{
		Name: html.EscapeString(name),
		Heroes: heroes,
		Cards: cards,
	}, nil
}

func readSerializedCard(bytes []byte, startIndex *int, endIndex int, prevCardBase *int) (Card, error) {
	if *startIndex > endIndex {
		return Card{}, errors.New("startIndex can not be greater than endIndex")
	}

	header := int(bytes[*startIndex])
	*startIndex++
	hasExtendedCount := (header >> 6) == 0x03

	cardDelta, err := readVarEncodedUint32(header, 5, bytes, startIndex, endIndex)

	if err != nil {
		return Card{}, errors.New("Can't read card delta")
	}

	outCardID := *prevCardBase + int(cardDelta)

	outCount := (header >> 6) + 1
	if hasExtendedCount {
		outCountExtended, err := readVarEncodedUint32(0, 0, bytes, startIndex, endIndex)
		if err != nil {
			return Card{}, errors.New("Can't read extended count")
		}

		outCount = int(outCountExtended)
	}


	*prevCardBase = outCardID

	return Card {
		ID: outCardID,
		Count: outCount,
	}, nil
}