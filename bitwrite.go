package artifactdeck

import "errors"

func extractNBitsWithCarry(value int, numBits uint) int {

	unLimitBit := 1 << numBits
	unResult := (value & (unLimitBit - 1))

	if (value >= unLimitBit) {
		unResult |= unLimitBit
	}

	return unResult
}


func addByte(bytes []byte, intByte int) (bool, []byte) {

	if(intByte > 255) {
		return false, nil
	}

	bytes = append(bytes, byte(intByte))

	return true, bytes
}

func addRemainingNumberToBuffer(bytes []byte, value int, alreadyWrittenBits uint) ([]byte, error) {

	value >>= alreadyWrittenBits
	var numBytes uint

	for value > 0 {
		nextByte := extractNBitsWithCarry(value, 7)
		value >>= 7

		_, bytes = addByte(bytes, nextByte)

		numBytes++
	}

	return bytes, nil
}


func addCardToBuffer(bytes []byte, count int, value int) ([]byte, error) {
	if count == 0 {
		return nil, errors.New("Invalid count")
	}

	countBytesStart := len(bytes)

	firstByteMaxCount := 0x03
	extendedCount := (count - 1) >= firstByteMaxCount

	firstByteCount := count -1
	if extendedCount {
		firstByteCount = firstByteMaxCount
	}

	firstByte := firstByteCount << 6
	firstByte |= extractNBitsWithCarry(value, 5)


	ok, bytes := addByte(bytes, firstByte)
	if !ok {
		return nil, errors.New("Can't append firstByte")
	}

	bytes, err := addRemainingNumberToBuffer(bytes, value, 5)
	if err != nil {
	 return nil, errors.New("Can't append rest of the number with a carry flag")
	}

	if extendedCount {
		bytes, err = addRemainingNumberToBuffer(bytes, count, 0)
		if err != nil {
		 return nil, errors.New("Can't append the remaining count")
		}
	}

	countBytesEnd := len(bytes)

	if countBytesEnd - countBytesStart > 11 {
		return nil, errors.New("Something went horribly wrong")
	}

	return bytes, nil
}

func computeChecksum(bytes []byte, numBytes int) int {
	checksum := 0
	for addCheck := headerSize; addCheck < numBytes + headerSize; addCheck++ {
		byte := bytes[addCheck]
		checksum += int(byte)
	}

	return checksum
}