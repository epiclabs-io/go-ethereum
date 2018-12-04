package flare

import "github.com/ethereum/go-ethereum/swarm/storage"

type HandlerParams struct {
	chunkStore *storage.NetStore
}

type Handler interface {
}

type handler struct {
	HandlerParams
}

func NewHandler(params *HandlerParams) Handler {
	return &handler{
		HandlerParams: *params,
	}

}

func (h *handler) Validate(chunkAddr storage.Address, data []byte) bool {

	var id ID
	err := id.UnmarshalBinary(chunkAddr)
	if err != nil {
		// warn
		return false
	}

	var v Value
	err = v.UnmarshalBinary(data)
	if err != nil {
		// warn
		return false
	}

	return v.Verify(&id)
}
