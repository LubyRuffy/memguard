package memguard

// import (
// 	"os"
// 	"unsafe"

// 	"github.com/awnumar/memguard/core"
// 	"github.com/awnumar/memguard/crypto"
// 	"github.com/awnumar/memguard/memcall"
// )

// /*
// NewEnclave creates a sealed and encrypted container from a given LockedBuffer. The given LockedBuffer will be automatically destroyed.

// Enclaves are immutable. To change the contents of an Enclave, call the View method to get a LockedBuffer with the contents and then create a new Enclave after making modifications.
// */
// func NewEnclave(b LockedBuffer) (Enclave, error) {
// 	b.Lock()                       // Grab the lock
// 	e, err := newEnclave(b.buffer) // Will error if LockedBuffer is destroyed
// 	if err != nil {
// 		return nil, ErrDestroyed
// 	}
// 	b.Unlock()    // Release the lock
// 	b.Destroy()   // Destroy the given LockedBuffer
// 	return e, nil // Return the shiny new Enclave
// }

// /*
// NewEnclaveFromBytes creates a sealed and encrypted container containing the given data. The original slice is wiped afterwards.

// If the length of the given buffer is zero, the call will return an ErrInvalidLength.
// */
// func NewEnclaveFromBytes(buf []byte) (Enclave, error) {
// 	return newEnclave(buf) // Will error if given a zero-length buffer
// }

// /*
// View decrypts the contents of an Enclave and returns the data inside a secure, mutable, LockedBuffer. You should destroy the returned LockedBuffer using the Destroy method as soon as you no longer need access to the contents.
// */
// func (b Enclave) View() *LockedBuffer {
// 	return newLockedBufferFromEnclave(b)
// }

// /*
// NewImmutable creates a new LockedBuffer of a specified size. The created object will be immutable, and this state can be toggled with the MakeImmutable and MakeMutable methods.

// If the given length is less than one, the call will return an ErrInvalidLength.
// */
// func NewImmutable(size int) (*LockedBuffer, error) {
// 	b, err := newLockedBuffer(size) // Will error if size < 1
// 	if err != nil {
// 		return nil, err
// 	}
// 	b.MakeImmutable() // Make the memory immutable
// 	return b, nil     // Return the LockedBuffer
// }

// /*
// NewMutable creates a new Enclave of a specified size. The created object will be mutable and sealed, and these states can be toggled with the MakeImmutable, MakeMutable, Unseal, and Reseal methods.

// If the given length is less than one, the call will return an ErrInvalidLength.
// */
// func NewMutable(size int) (LockedBuffer, error) {
// 	// Create a new Enclave.
// 	b, err := newContainer(size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Return the Enclave object.
// 	return b, nil
// }

// /*
// NewImmutableFromBytes is identical to NewImmutable but for the fact that the created Enclave is of the same length and has the same contents as a given slice. The slice is wiped after the bytes have been copied over.

// If the size of the slice is zero, the call will return an ErrInvalidLength.
// */
// func NewImmutableFromBytes(buf []byte) (LockedBuffer, error) {
// 	// Create a new Enclave.
// 	b, err := newContainer(len(buf))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Copy the bytes from buf, wiping afterwards.
// 	b.Move(buf)

// 	// Seal the Enclave.
// 	b.reseal()

// 	// Make the memory immutable.
// 	b.MakeImmutable()

// 	// Return a pointer to the Enclave.
// 	return b, nil
// }

// /*
// NewMutableFromBytes is identical to NewMutable but for the fact that the created Enclave is of the same length and has the same contents as a given slice. The slice is wiped after the bytes have been copied over.

// If the size of the slice is zero, the call will return an ErrInvalidLength.
// */
// func NewMutableFromBytes(buf []byte) (LockedBuffer, error) {
// 	// Create a new Enclave.
// 	b, err := newContainer(len(buf))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Copy the bytes from buf, wiping afterwards.
// 	b.Move(buf)

// 	// Seal the Enclave.
// 	b.reseal()

// 	// Return a pointer to the Enclave.
// 	return b, nil
// }

// /*
// NewImmutableRandom is identical to NewImmutable except the created Enclave is filled with cryptographically-secure pseudo-random bytes instead of zeroes.

// If the given length is less than one, the call will return an ErrInvalidLength.
// */
// func NewImmutableRandom(size int) (LockedBuffer, error) {
// 	// Create a new Enclave for the key.
// 	b, err := newContainer(size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fill it with random data.
// 	if err := crypto.MemScr(b.plaintext); err != nil {
// 		SafePanic(err)
// 	}

// 	// Seal the Enclave.
// 	b.reseal()

// 	// Make the memory immutable.
// 	b.MakeImmutable()

// 	// Return the Enclave.
// 	return b, nil
// }

// /*
// NewMutableRandom is identical to NewMutable except the created Enclave is filled with cryptographically-secure pseudo-random bytes instead of zeroes.

// If the given length is less than one, the call will return an ErrInvalidLength.
// */
// func NewMutableRandom(size int) (LockedBuffer, error) {
// 	// Create a new Enclave for the key.
// 	b, err := newContainer(size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fill it with random data.
// 	if err := crypto.MemScr(b.plaintext); err != nil {
// 		SafePanic(err)
// 	}

// 	// Seal the Enclave.
// 	b.reseal()

// 	// Return the Enclave.
// 	return b, nil
// }

// /*
// Unseal decrypts and unseals a given Enclave.

// All new Enclaves are sealed by default and you should use this method to temporarily unseal them, so as to view or modify their contents, and then call Reseal again as soon as possible. There is no need to call Unseal if you are using MemGuard's own API since internal functions know how to handle sealed containers and will unseal and reseal them for you automatically. Instead, call Unseal before passing the contents of a container to another API, and reiterating, call Reseal again soon after.

// Calling Unseal on an unsealed Enclave will return ErrUnsealed. This is because the developer has assumed that a container is sealed when it isn't, and that is a dangerous assumption.

// If the Enclave is immutable, Unseal will automatically work around it and preserve the immutability state.
// */
// func (b LockedBuffer) Unseal() error {
// 	// Attain the mutex.
// 	b.Lock()
// 	defer b.Unlock()

// 	return b.unseal()
// }

// /*
// Reseal re-encrypts and reseals a given Enclave.

// All new Enclaves are sealed by default and you should use Unseal to unseal them them, so as to view or modify their contents, and then call this method again as soon as possible.

// Calling Reseal on a sealed Enclave does nothing. If the given Enclave is immutable, Reseal will automatically work around it and preserve the immutability state.
// */
// func (b LockedBuffer) Reseal() error {
// 	// Attain the mutex.
// 	b.Lock()
// 	defer b.Unlock()

// 	return b.reseal()
// }

// /*
// Bytes returns a slice that references the secure, protected portion of memory.

// If the Enclave that you call Bytes on is sealed, the data returned will be random and so Unseal should be called first (promptly followed by Reseal when done). Otherwise, if the Enclave has been destroyed then the returned slice will be nil (it will have a length and capacity of zero).

// If a function that you're using requires an array, you can cast the slice to an array and then pass around a pointer:

//     // Make sure the size of the array matches the size of the buffer.
//     // In this example that size is 16. This is VERY important.
// 	keyArrayPtr := (*[16]byte)(unsafe.Pointer(&b.Bytes()[0]))

// 	// Pass around the array pointer WITHOUT dereferencing it.
// 	Encrypt(plaintext, keyArrayPtr)

// Make sure that you do not dereference the pointer and pass around the resulting value as this will leave copies all over the place.
// */
// func (b LockedBuffer) Bytes() []byte {
// 	return b.plaintext
// }

// /*
// Uint16 returns a slice (of type []uint16) that references the secure, protected portion of memory.

// The Enclave must be a multiple of 2 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Uint16() ([]uint16, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%2 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 2, b.Size() / 2}

// 	// Return the new slice.
// 	return *(*[]uint16)(unsafe.Pointer(&sl)), nil
// }

// /*
// Uint32 returns a slice (of type []uint32) that references the secure, protected portion of memory.

// The Enclave must be a multiple of 4 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Uint32() ([]uint32, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%4 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 4, b.Size() / 4}

// 	// Return the new slice.
// 	return *(*[]uint32)(unsafe.Pointer(&sl)), nil
// }

// /*
// Uint64 returns a slice (of type []uint64) that references the secure, protected portion of memory.

// The Enclave must be a multiple of 8 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Uint64() ([]uint64, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%8 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 8, b.Size() / 8}

// 	// Return the new slice.
// 	return *(*[]uint64)(unsafe.Pointer(&sl)), nil
// }

// /*
// Int8 returns a slice (of type []int8) that references the secure, protected portion of memory.
// */
// func (b LockedBuffer) Int8() ([]int8, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size(), b.Size()}

// 	// Return the new slice.
// 	return *(*[]int8)(unsafe.Pointer(&sl)), nil
// }

// /*
// Int16 returns a slice (of type []int16) that references the secure, protected portion of memory.

// The Enclave must be a multiple of 2 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Int16() ([]int16, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%2 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 2, b.Size() / 2}

// 	// Return the new slice.
// 	return *(*[]int16)(unsafe.Pointer(&sl)), nil
// }

// /*
// Int32 returns a slice (of type []int32) that references the secure, protected portion of memory.

// The Enclave must be a multiple of 4 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Int32() ([]int32, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%4 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 4, b.Size() / 4}

// 	// Return the new slice.
// 	return *(*[]int32)(unsafe.Pointer(&sl)), nil
// }

// /*
// Int64 returns a slice (of type []int64) that references the secure, protected portion of memory.

// The Enclave must also be a multiple of 8 bytes in length or an error will be returned.
// */
// func (b LockedBuffer) Int64() ([]int64, error) {
// 	// Attain the mutex lock.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Check to see if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check if it's sealed.
// 	if b.sealed {
// 		return nil, ErrSealed
// 	}

// 	// Check to see if it's an appropriate length.
// 	if len(b.plaintext)%8 != 0 {
// 		return nil, ErrInvalidConversion
// 	}

// 	// Perform the conversion.
// 	var sl = struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}{uintptr(unsafe.Pointer(&b.plaintext[0])), b.Size() / 8, b.Size() / 8}

// 	// Return the new slice.
// 	return *(*[]int64)(unsafe.Pointer(&sl)), nil
// }

// /*
// IsSealed returns a boolean value indicating if an Enclave is Sealed.
// */
// func (b LockedBuffer) IsSealed() bool {
// 	// Attain the mutex.
// 	b.RLock()
// 	defer b.RUnlock()

// 	return b.sealed
// }

// /*
// IsMutable returns a boolean value indicating if an Enclave is mutable.
// */
// func (b LockedBuffer) IsMutable() bool {
// 	// Get a mutex lock on this Enclave.
// 	b.RLock()
// 	defer b.RUnlock()

// 	return b.mutable
// }

// /*
// IsDestroyed returns a boolean value indicating if an Enclave has been destroyed.
// */
// func (b LockedBuffer) IsDestroyed() bool {
// 	// Get a mutex lock on this Enclave.
// 	b.RLock()
// 	defer b.RUnlock()

// 	// Return the appropriate value.
// 	return len(b.plaintext) == 0
// }

// /*
// EqualBytes compares an Enclave to a byte slice in constant time. If called on a sealed Enclave, EqualBytes will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) EqualBytes(buf []byte) (bool, error) {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return false, ErrDestroyed
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Do a time-constant comparison and return the result.
// 	return crypto.Equal(b.plaintext, buf), nil
// }

// /*
// MakeImmutable asks the kernel to mark the Enclave's contents immutable. Any subsequent attempts to modify this memory will result in the kernel raising an access violation and terminating the process. To make the buffer mutable again, call MakeMutable.

// The API will respect your mutability state and so if a MemGuard function that needs to modify an Enclave is handed one that is immutable, it will return ErrImmutable.
// */
// func (b LockedBuffer) MakeImmutable() error {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return ErrDestroyed
// 	}

// 	if b.mutable {
// 		// Mark the memory as mutable.
// 		if err := memcall.Protect(getAllMemory(b)[pageSize:pageSize+roundToPageSize(len(b.plaintext)+32)], true, false); err != nil {
// 			SafePanic(err)
// 		}

// 		// Tell everyone about the change we made.
// 		b.mutable = false
// 	}

// 	// Everything went well.
// 	return nil
// }

// /*
// MakeMutable asks the kernel to mark the Enclave's contents mutable. To make the buffer immutable again, call MakeImmutable.
// */
// func (b LockedBuffer) MakeMutable() error {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return ErrDestroyed
// 	}

// 	if !b.mutable {
// 		// Mark the memory as mutable.
// 		if err := memcall.Protect(getAllMemory(b)[pageSize:pageSize+roundToPageSize(len(b.plaintext)+32)], true, true); err != nil {
// 			SafePanic(err)
// 		}

// 		// Tell everyone about the change we made.
// 		b.mutable = true
// 	}

// 	// Everything went well.
// 	return nil
// }

// /*
// Copy copies bytes from a byte slice into an Enclave, in constant-time. Just like Golang's built-in copy function, Copy only copies up to the smallest of the two buffers.

// Copy does not wipe the original slice so using Move should be favoured unless you have a specific reason. You should aim to call WipeBytes on the original slice as soon as possible after copying.

// If the Enclave is marked as immutable, the call will fail and return an ErrImmutable. If the Enclave is sealed, Copy will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) Copy(buf []byte) error {
// 	// Just call CopyAt with a zero offset.
// 	return b.CopyAt(buf, 0)
// }

// /*
// CopyAt is identical to Copy but it copies into the Enclave at a specified offset.
// */
// func (b LockedBuffer) CopyAt(buf []byte, offset int) error {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return ErrDestroyed
// 	}

// 	// Check if it's immutable.
// 	if !b.mutable {
// 		return ErrImmutable
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Do a time-constant copying of the bytes.
// 	crypto.Copy(b.plaintext[offset:], buf)

// 	return nil
// }

// /*
// Move moves bytes from a byte slice into an Enclave in constant-time. Just like Golang's built-in copy function, Move only moves up to the smallest of the two buffers.

// Unlike Copy, Move wipes the entire original slice after copying, and so it should be favoured unless you have a specific reason.

// If the Enclave is marked as immutable, the call will fail and return an ErrImmutable. If the Enclave is sealed, Move will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) Move(buf []byte) error {
// 	// Just call MoveAt with a zero offset.
// 	return b.MoveAt(buf, 0)
// }

// /*
// MoveAt is identical to Move but it copies into the Enclave at a specified offset.
// */
// func (b LockedBuffer) MoveAt(buf []byte, offset int) error {
// 	// Copy buf into the Enclave.
// 	if err := b.CopyAt(buf, offset); err != nil {
// 		return err
// 	}

// 	// Wipe the old bytes.
// 	crypto.MemClr(buf)

// 	// Everything went well.
// 	return nil
// }

// /*
// FillRandomBytes fills an Enclave with cryptographically-secure pseudo-random bytes. If the Enclave is sealed, it will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) FillRandomBytes() error {
// 	// Just call FillRandomBytesAt.
// 	return b.FillRandomBytesAt(0, b.Size())
// }

// /*
// FillRandomBytesAt fills an Enclave with cryptographically-secure pseudo-random bytes, starting at an offset and ending after a given number of bytes. If the Enclave is sealed, it will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) FillRandomBytesAt(offset, length int) error {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return ErrDestroyed
// 	}

// 	// Check if it's immutable.
// 	if !b.mutable {
// 		return ErrImmutable
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Fill with random bytes.
// 	if err := crypto.MemScr(b.plaintext[offset : offset+length]); err != nil {
// 		SafePanic(err)
// 	}

// 	// Everything went well.
// 	return nil
// }

// /*
// Destroy performs some security checks, securely wipes the contents of, and then releases a LockedBuffer's memory back to the OS. If a security check fails, the process will attempt to wipe all it can before safely panicking.

// This function should be called on all LockedBuffers before exiting. DestroyAll is designed for this purpose, as is CatchInterrupt, SafePanic, and SafeExit. We recommend using a mixture of them depending on your situation.

// If the Enclave has already been destroyed, subsequent calls are idempotent.
// */
// func (b LockedBuffer) Destroy() {
// 	b.raw.Destroy()
// }

// /*
// Size returns an integer representing the total length, in bytes, of an Enclave.

// If this size is zero, it is safe to assume that the Enclave has been destroyed.
// */
// func (b LockedBuffer) Size() int {
// 	return len(b.plaintext)
// }

// /*
// Wipe wipes an Enclave's contents by overwriting the buffer with zeroes. If the Enclave is sealed, it will automatically unseal and reseal it for you.
// */
// func (b LockedBuffer) Wipe() error {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return ErrDestroyed
// 	}

// 	// Check if it's immutable.
// 	if !b.mutable {
// 		return ErrImmutable
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Wipe the buffer.
// 	crypto.MemClr(b.plaintext)

// 	// Everything went well.
// 	return nil
// }

// /*
// Concatenate takes two Enclaves, concatenates them, and returns a sealed container. The original Enclaves are preserved.

// If one of the given Enclaves is immutable, the resulting Enclave will also be immutable. If an Enclave is sealed, Concatenate will automatically unseal and reseal it for you.
// */
// func Concatenate(a, b LockedBuffer) (LockedBuffer, error) {
// 	// Get a mutex lock on the Enclaves.
// 	a.Lock()
// 	b.Lock()
// 	defer a.Unlock()
// 	defer b.Unlock()

// 	// Check if either are destroyed.
// 	if len(a.plaintext) == 0 || len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check to see if either are sealed.
// 	if a.sealed {
// 		a.unseal()
// 		defer a.reseal()
// 	}
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Create a new Enclave to hold the concatenated value.
// 	c, _ := newContainer(len(a.plaintext) + len(b.plaintext))

// 	// Copy the values across.
// 	c.Copy(a.plaintext)
// 	c.CopyAt(b.plaintext, len(a.plaintext))

// 	// Seal the container.
// 	c.reseal()

// 	// Set permissions accordingly.
// 	if !a.mutable || !b.mutable {
// 		c.MakeImmutable()
// 	}

// 	// Return the resulting Enclave.
// 	return c, nil
// }

// /*
// Duplicate takes an Enclave and creates a new one with the same contents and mutability state as the original.

// If the Enclave is sealed, it will automatically unseal and reseal it for you. The returned Enclave will be sealed regardless.
// */
// func Duplicate(b LockedBuffer) (LockedBuffer, error) {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Create new Enclave.
// 	newBuf, _ := newContainer(b.Size())

// 	// Copy bytes into it.
// 	newBuf.Copy(b.plaintext)

// 	// Seal it.
// 	newBuf.reseal()

// 	// Set permissions accordingly.
// 	if !b.mutable {
// 		newBuf.MakeImmutable()
// 	}

// 	// Return duplicated.
// 	return newBuf, nil
// }

// /*
// Equal compares the contents of two Enclaves in constant time. If either are sealed, it will automatically unseal and reseal them for you.
// */
// func Equal(a, b LockedBuffer) (bool, error) {
// 	// Get a mutex lock on the Enclaves.
// 	a.Lock()
// 	b.Lock()
// 	defer a.Unlock()
// 	defer b.Unlock()

// 	// Check if either are destroyed.
// 	if len(a.plaintext) == 0 || len(b.plaintext) == 0 {
// 		return false, ErrDestroyed
// 	}

// 	// Check to see if either are sealed.
// 	if a.sealed {
// 		a.unseal()
// 		defer a.reseal()
// 	}
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Do a time-constant comparison on the two buffers; return the result.
// 	return crypto.Equal(a.plaintext, b.plaintext), nil
// }

// /*
// Split takes an Enclave, splits it at a specified offset, and then returns the two newly created Enclaves. The mutability state of the original is preserved in the new (sealed) Enclaves, and the original Enclave is not destroyed. If the given Enclave is sealed, Split will automatically unseal and reseal it for you.
// */
// func Split(b LockedBuffer, offset int) (LockedBuffer, LockedBuffer, error) {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, nil, ErrDestroyed
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Create two new Enclaves.
// 	firstBuf, err := newContainer(len(b.plaintext[:offset]))
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	secondBuf, err := newContainer(len(b.plaintext[offset:]))
// 	if err != nil {
// 		firstBuf.Destroy()
// 		return nil, nil, err
// 	}

// 	// Copy the values into them.
// 	firstBuf.Copy(b.plaintext[:offset])
// 	secondBuf.Copy(b.plaintext[offset:])

// 	// Seal them.
// 	firstBuf.reseal()
// 	secondBuf.reseal()

// 	// Copy over permissions.
// 	if !b.mutable {
// 		firstBuf.MakeImmutable()
// 		secondBuf.MakeImmutable()
// 	}

// 	// Return the new Enclaves.
// 	return firstBuf, secondBuf, nil
// }

// /*
// Grow takes an Enclave and returns a new one that's n bytes larger. The contents of the given Enclave would be preserved in the first b.Size() bytes of the newly created one.

// The mutability state of the original is preserved in the new (sealed) Enclave, and the original Enclave is not destroyed. If Trim is called on a sealed Enclave, it will automatically unseal and reseal it for you.
// */
// func Grow(b LockedBuffer, n int) (LockedBuffer, error) {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Create new Enclave and copy over the old.
// 	newBuf, err := newContainer(len(b.plaintext) + n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	newBuf.Copy(b.plaintext)

// 	// Seal it up.
// 	newBuf.reseal()

// 	// Copy over permissions.
// 	if !b.mutable {
// 		newBuf.MakeImmutable()
// 	}

// 	// Return the new Enclave.
// 	return newBuf, nil
// }

// /*
// Trim shortens an Enclave according to the given specifications. It takes an offset and a size as arguments. The resulting Enclave starts at index [offset] and ends at index [offset+size].

// The mutability state of the original is preserved in the new (sealed) Enclave, and the original Enclave is not destroyed. If Trim is called on a sealed Enclave, it will automatically unseal and reseal it for you.
// */
// func Trim(b LockedBuffer, offset, size int) (LockedBuffer, error) {
// 	// Get a mutex lock on this Enclave.
// 	b.Lock()
// 	defer b.Unlock()

// 	// Check if it's destroyed.
// 	if len(b.plaintext) == 0 {
// 		return nil, ErrDestroyed
// 	}

// 	// Check to see if it's sealed.
// 	if b.sealed {
// 		b.unseal()
// 		defer b.reseal()
// 	}

// 	// Create new Enclave and copy over the old.
// 	newBuf, err := newContainer(size)
// 	if err != nil {
// 		return nil, err
// 	}
// 	newBuf.Copy(b.plaintext[offset : offset+size])

// 	// Seal it up.
// 	newBuf.reseal()

// 	// Copy over permissions.
// 	if !b.mutable {
// 		newBuf.MakeImmutable()
// 	}

// 	// Return the new Enclave.
// 	return newBuf, nil
// }

// /*
// WipeBytes zeroes out a given buffer. It is recommended that you call WipeBytes on slices after utilizing the Copy or CopyAt methods.

// Due to the nature of memory allocated by the Go runtime, WipeBytes cannot guarantee that the data does not exist elsewhere in memory. Therefore, your program should aim to (as far as possible) store sensitive data only in Enclaves, which expose their own Wipe method.
// */
// func WipeBytes(b []byte) {
// 	crypto.MemClr(b)
// }

// /*
// DestroyAll calls Destroy on all Enclaves that have not already been destroyed.

// CatchInterrupt and SafeExit both call DestroyAll before exiting.
// */
// func DestroyAll() {
// 	// Get a Mutex lock on enclaves, and get a copy.
// 	enclavesMutex.RLock()
// 	containers := make([]LockedBuffer, len(enclaves))
// 	copy(containers, enclaves)
// 	enclavesMutex.RUnlock()

// 	// Destroy them all.
// 	for _, b := range containers {
// 		b.Destroy()
// 	}
// }

// /*
// CatchInterrupt allows you to specify a function to be executed in the event that the process receives an interrupt signal. The thread will call SafeExit immediately after executing the given function.
// */
// func CatchInterrupt(f func()) {
// 	interruptHandler = f
// }

// /*
// SafePanic is identical to Go's panic except it cleans up all it can before calling panic.
// */
// func SafePanic(v interface{}) {
// 	core.Panic(v)
// }

// /*
// SafeExit exits the program with a specified exit-code, but cleans up first.
// */
// func SafeExit(c int) {
// 	// Clean-up protected memory.
// 	DestroyAll()

// 	// Destroy the global protection values.
// 	subclaves.canary.destroy()
// 	subclaves.enckey.destroy()

// 	// Exit with a specified exit-code.
// 	os.Exit(c)
// }

// /*
// DisableUnixCoreDumps disables core-dumps.

// Since core-dumps are only relevant on Unix systems, if DisableUnixCoreDumps is called on any other system it will do nothing and return immediately.

// This function is precautonary as core-dumps are usually disabled by default on most systems.
// */
// func DisableUnixCoreDumps() error {
// 	return memcall.DisableCoreDumps()
// }
