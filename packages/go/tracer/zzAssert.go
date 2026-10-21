package tracer

// Compile-time assertions that the concrete types satisfy the interfaces the
// types module declares. Without these, a method rename on either side drifts
// silently: the package still builds, and only a caller that happens to assign
// across the seam notices.
var (
	_ TracerLogInterface       = (*TracerLog)(nil)
	_ TracerLogFilterInterface = (*TracerLogFilter)(nil)
	_ ReadableTracerLog        = (*Tracer)(nil)
	_ TracerLogSource          = (*FileSource)(nil)
)
