package types

// HotReloadable is the dispatch seam: a producer calls Dispatch, the consumer
// implements OnDispatch. It is deliberately not HMR-specific — any message can
// travel the same path.
type HotReloadable interface {
	Dispatch(message string)
	OnDispatch(message string)
}

// HmrRefreshMessage asks a connected browser client to reload the page.
const HmrRefreshMessage = "riot::refresh"
