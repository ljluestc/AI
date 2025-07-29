package vision

// CursorState represents the state of a cursor in the vision module.
type CursorState struct {
	// Position in 3D space
	Position [3]float64
	// Visibility flag (matches .Visible usage)
	Visible bool
	// Optional: mask data for compatibility with tests
	Mask [][][]float64
}

// IsVisible returns the visibility status of the cursor.
func (c *CursorState) IsVisible() bool {
	return c.Visible
}

// CursorStateEnhanced represents the state of the cursor with additional attributes.
type CursorStateEnhanced struct {
	Position     [3]float64    // 3D position
	Velocity     [3]float64    // Optional: velocity vector
	Acceleration [3]float64    // Optional: acceleration vector
	Confidence   float64       // Detection confidence
	Mask         [][][]float64 // Optional: mask data
}

// To allow conversion from CursorStateEnhanced to CursorState,
// add a helper function.
func (e *CursorStateEnhanced) ToCursorState() *CursorState {
	return &CursorState{
		Position: e.Position,
		Visible:  false, // fallback: CursorStateEnhanced does not track visibility
		Mask:     e.Mask,
	}
}
