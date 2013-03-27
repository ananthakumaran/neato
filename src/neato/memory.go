package main

type readCallbackFunc func(address uint16) byte
type writeCallbackFunc func(address uint16, val byte)

type Memory struct {
	store          []byte
	readCallbacks  []*readCallbackFunc
	writeCallbacks []*writeCallbackFunc
}

func NewMemory(size uint16) *Memory {
	memory := Memory{}
	memory.store = make([]byte, uint32(size)+1)
	memory.readCallbacks = make([]*readCallbackFunc, uint32(size)+1)
	memory.writeCallbacks = make([]*writeCallbackFunc, uint32(size)+1)
	return &memory
}

func (memory *Memory) Read(address uint16) byte {
	if callback := memory.readCallbacks[address]; callback != nil {
		return (*callback)(address)
	}
	return memory.store[address]
}

func (memory *Memory) Write(address uint16, value uint8) {
	if callback := memory.writeCallbacks[address]; callback != nil {
		(*callback)(address, value)
	} else {
		memory.store[address] = value
	}
}

func (memory *Memory) Copy(start, end int, from []byte) {
	copy(memory.store[start:end], from)
}

func (memory *Memory) ReadCallback(start, end uint16, callback readCallbackFunc) {

	for i := uint32(start); i <= uint32(end); i++ {
		memory.readCallbacks[i] = &callback
	}
}

func (memory *Memory) WriteCallback(start, end uint16, callback writeCallbackFunc) {
	for i := uint32(start); i <= uint32(end); i++ {
		memory.writeCallbacks[i] = &callback
	}
}

func (memory *Memory) Mirror(start, end, mstart, mend uint16) {
	interval := uint32(end - start)
	tempStart := uint32(mstart)

	for ; tempStart+interval <= uint32(mend); tempStart += (interval + 1) {
		tempMirrorStart := uint16(tempStart)
		tempEnd := tempStart + interval

		memory.ReadCallback(tempMirrorStart, uint16(tempEnd), func(address uint16) byte {
			//info("mirror read Original %04X destination %04X\n", address, start+address-tempMirrorStart)
			return memory.Read(start + address - tempMirrorStart)
		})

		memory.WriteCallback(tempMirrorStart, uint16(tempEnd), (func(address uint16, val byte) {
			memory.Write(start+address-tempMirrorStart, val)
		}))
	}

	if uint16(tempStart-1) != mend {
		debug("mstart %x mend %x", tempStart-1, mend)
		fatal("invalid mirror arguments")
	}
}
