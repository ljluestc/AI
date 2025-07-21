package rl

import (
	"errors"
	"fmt"
)

// BuddyMemory implements a buddy memory allocation system
type BuddyMemory struct {
	Size      int
	Allocated map[int]bool
}

// CreateBuddyMemory creates a new BuddyMemory with the specified size
func CreateBuddyMemory(size int) *BuddyMemory {
	return &BuddyMemory{
		Size:      size,
		Allocated: make(map[int]bool),
	}
}

// Allocate allocates a block of memory of the given size
func (bm *BuddyMemory) Allocate(size int) (int, error) {
	if size <= 0 {
		return -1, errors.New("size must be positive")
	}
	if size > bm.Size {
		return -1, fmt.Errorf("requested size %d exceeds available memory %d", size, bm.Size)
	}

	// Simple implementation: find the first free block
	for addr := 0; addr <= bm.Size-size; addr += size {
		if !bm.isAllocated(addr, addr+size) {
			// Mark this block as allocated
			for i := addr; i < addr+size; i++ {
				bm.Allocated[i] = true
			}
			return addr, nil
		}
	}

	return -1, errors.New("no available memory block of requested size")
}

// Free frees the memory block at the given address
func (bm *BuddyMemory) Free(addr int) error {
	if addr < 0 || addr >= bm.Size {
		return fmt.Errorf("invalid address: %d", addr)
	}

	if !bm.Allocated[addr] {
		return fmt.Errorf("memory at address %d is not allocated", addr)
	}

	// Find the size of the allocated block
	size := 0
	for i := addr; i < bm.Size && bm.Allocated[i]; i++ {
		size++
		bm.Allocated[i] = false
	}

	return nil
}

// isAllocated checks if any part of the memory range is allocated
func (bm *BuddyMemory) isAllocated(start, end int) bool {
	for i := start; i < end; i++ {
		if bm.Allocated[i] {
			return true
		}
	}
	return false
}
