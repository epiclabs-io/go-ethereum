package flare

import "github.com/ethereum/go-ethereum/swarm/storage"

type HandlerParams struct {
	ChunkStore *storage.NetStore
}

type Handler interface {
	storage.ChunkValidator
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

	var v Value
	err := v.UnmarshalBinary(data)
	if err != nil {
		// warn
		return false
	}

	return v.Verify(chunkAddr)
}

func (h *handler) Put() error {
	return nil
}

func (h *handler) Get(addr Address) error {

	return nil
}
