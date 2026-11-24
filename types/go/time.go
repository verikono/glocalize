package types

// a period of time measured in seconds
type Seconds int

// a period of time measured in miliseconds
type MilliSeconds int

type NanoSeconds int64

// a position in time, measured by the number of seconds since the start of Jan 1st 1970
type EpochSeconds Seconds

// a position in time, measured by the number of seconds since the start of Jan 1st 1970
type EpochMiliSeconds MilliSeconds

// a position in time, measured by the number of seconds since the start of Jan 1st 1970
type EpochNanoSeconds NanoSeconds
