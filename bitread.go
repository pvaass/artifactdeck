package artifactdeck

import "errors"

func readVarEncodedUint32(baseValue int, baseBits int, data []byte, startIndex *int, endIndex int) (int, error) {
	outValue := 0

	deltaShift := 0

	if (baseBits == 0) || readBitsChunk(baseValue, baseBits, deltaShift, &outValue) {
		deltaShift += baseBits

		for true {
			if *startIndex > endIndex {
				return 0, errors.New("No more room")
			}

			nextByte := int(data[*startIndex])
			*startIndex++

			if !readBitsChunk(nextByte, 7, deltaShift, &outValue) {
				break
			}

			deltaShift += 7
		}

	}


	return outValue, nil
}

func readBitsChunk(chunk int, numBits int, currShift int, outBits *int) bool {
	continueBit := 1 << uint(numBits)
	newBits := chunk & (continueBit -1)
	
	*outBits |= (newBits << uint(currShift))

	return (chunk & continueBit) != 0
}