package utility



type ReservationError string

func (e ReservationError) Error() string {
	return string(e)
}

const (
	ErrReservationNotFound    = ReservationError("Failed to find reservation")
	ErrCancellationTooSoon    = ReservationError("Cancellation allowed only for events more than 24 hours in advance")
	ErrFailedToDelete         = ReservationError("Failed to delete reservation")
	ErrFailedToUpdateSeat     = ReservationError("Failed to update seat reservation status")
	ErrFailedToCommit         = ReservationError("Failed to commit transaction")
)
