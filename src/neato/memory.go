package main

type readCallbackFunc func(address uint16) byte
type readCallback struct {
	start uint16
	end   uint16
	fn    readCallbackFunc
}

type writeCallbackFunc func(address uint16, val byte)
type writeCallback struct {
	start uint16
	end   uint16
	fn    writeCallbackFunc
}

type Memory struct {
	store          []byte
	readCallbacks  []readCallback
	writeCallbacks []writeCallback
}

func newMemory(size uint16) *Memory {
	memory := Memory{}
	memory.store = make([]byte, uint32(size)+1)
	return &memory
}

func (memory *Memory) read(address uint16) byte {
	for _, callback := range memory.readCallbacks {
		if callback.start <= address && callback.end >= address {
			return callback.fn(address)
		}
	}

	return memory.store[address]
}

func (memory *Memory) write(address uint16, value uint8) {
	for _, callback := range memory.writeCallbacks {
		if callback.start <= address && callback.end >= address {
			callback.fn(address, value)
			return
		}
	}
	memory.store[address] = value
}

func (memory *Memory) copy(start, end int, from []byte) {
	copy(memory.store[start:end], from)
}

func (memory *Memory) readCallback(start, end uint16, callback readCallbackFunc) {
	memory.readCallbacks = append(memory.readCallbacks, readCallback{start, end, callback})
}

func (memory *Memory) writeCallback(start, end uint16, callback writeCallbackFunc) {
	memory.writeCallbacks = append(memory.writeCallbacks, writeCallback{start, end, callback})
}

func (memory *Memory) mirror(start, end, mstart, mend uint16) {
	interval := uint32(end - start)
	tempStart := uint32(mstart)

	for ; tempStart+interval <= uint32(mend); tempStart += (interval + 1) {
		tempEnd := tempStart + interval

		memory.readCallback(mstart, uint16(tempEnd), func(address uint16) byte {
			return memory.read(start + address - mstart)
		})

		memory.writeCallback(mstart, uint16(tempEnd), func(address uint16, val byte) {
			memory.write(start+address-mstart, val)
		})
	}

	if uint16(tempStart-1) != mend {
		debug("mstart %x mend %x", tempStart-1, mend)
		fatal("invalid mirror arguments")
	}

}
